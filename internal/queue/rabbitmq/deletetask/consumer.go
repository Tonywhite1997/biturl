package queue

import (
	"biturl/internal/helper"
	"biturl/internal/queue/rabbitmq"
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartDeleteWorker(conn *amqp.Connection, queueName string, deleteFunc func(ctx context.Context, shortcode string) error) {

	ch, err := conn.Channel()
	helper.FailOnError(err, "failed to open channel in insert worker")

	_, err = ch.QueueDeclare(rabbitmq.DeleteRedisQueueKey, true, false, false, false, nil)
	helper.FailOnError(err, "failed to declare queue in insert worker")

	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	helper.FailOnError(err, "failed to consume")

	defer helper.RecoverWorker()

	for msg := range msgs {
		var task DeleteTask
		if err := json.Unmarshal(msg.Body, &task); err != nil {
			fmt.Println("invalid message:", err)
			msg.Ack(false) // acknowledge to skip bad message
			continue
		}

		ctx := context.Background()
		var lastErr error

		for i := 0; i < 3; i++ {
			if err := deleteFunc(ctx, task.ShortCode); err == nil {
				fmt.Println("deleted redis key", task.ShortCode)
				msg.Ack(false)
				lastErr = nil
				break
			} else {
				lastErr = err
				fmt.Println("retry failed, attempt", i+1, "error:", err)
				time.Sleep(helper.RetryInterval(i))
			}
		}

		if lastErr != nil {
			fmt.Println("all retries failed for", lastErr)
			msg.Nack(false, true)
		}
	}

}
