package worker

import (
	"biturl/internal/queue/rabbitmq/insertstat"
	"biturl/internal/repository"
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartInsertWorker(conn *amqp.Connection, queueName string, insertFunc func(ctx context.Context, stats repository.Stats) error) {
	insertstat.StartInsertWorker(conn, queueName, insertFunc)
}
