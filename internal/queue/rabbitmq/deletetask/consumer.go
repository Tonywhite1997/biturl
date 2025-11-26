package queue

import (
	"biturl/internal/helper"
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartDeleteWorker(ch *amqp.Channel, queueName string, deleteFunc func(ctx context.Context, shortcode string) error) {

	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	helper.FailOnError(err, "failed to consume")

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("worker panicked:", r)
			}
		}()

		for msg := range msgs {
			var task DeleteTask
			if err := json.Unmarshal(msg.Body, &task); err != nil {
				fmt.Println("invalid message:", err)
				msg.Ack(false) // acknowledge to skip bad message
				continue
			}

			ctx := context.Background()
			success := false

			for i := 0; i < 3; i++ {
				if err := deleteFunc(ctx, task.ShortCode); err == nil {
					fmt.Println("deleted redis key", task.ShortCode)
					msg.Ack(false)
					success = true
					break
				} else {
					fmt.Println("retry failed, attempt", i+1, "error:", err)
					time.Sleep(5 * time.Second)
				}
			}

			if !success {
				fmt.Println("all retries failed for", task.ShortCode)
				msg.Nack(false, true) // requeue message
			}
		}
	}()
}
