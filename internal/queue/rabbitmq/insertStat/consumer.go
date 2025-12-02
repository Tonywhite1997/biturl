package insertstat

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"biturl/internal/helper"
	"biturl/internal/repository"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartInsertWorker(conn *amqp.Connection, queuename string, insertFunc func(ctx context.Context, stats repository.Stats) error) {

	ch, err := conn.Channel()
	helper.FailOnError(err, "failed to open channel for insert worker")

	_, err = ch.QueueDeclare(queuename, true, false, false, false, nil)
	helper.FailOnError(err, "failed to declare queue for insert worker")

	msgs, err := ch.Consume(
		queuename,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	helper.FailOnError(err, "could not consume queue")

	defer helper.RecoverWorker()

	for msg := range msgs {
		var stat repository.Stats

		if err := json.Unmarshal(msg.Body, &stat); err != nil {
			fmt.Println("invalid message:", err)
			msg.Ack(false)
			continue
		}

		ctx := context.Background()
		var lastErr error

		// Retry 3 times
		for i := 0; i < 3; i++ {
			if err := insertFunc(ctx, stat); err == nil {
				fmt.Println("stat inserted:", stat)
				msg.Ack(false)
				lastErr = nil
				break
			} else {
				lastErr = err
				fmt.Println("failed to insert, attempt", i+1, err)
				time.Sleep(helper.RetryInterval(i))
			}
		}

		if lastErr != nil {
			fmt.Println("all retries failed for:", stat, "error:", lastErr)
			msg.Nack(false, true)
		}
	}
}
