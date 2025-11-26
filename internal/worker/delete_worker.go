package worker

import (
	queue "biturl/internal/queue/rabbitmq/deletetask"
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartDeleteWorker(conn *amqp.Channel, queueName string, deleteFunc func(ctx context.Context, shortCode string) error) {
	queue.StartDeleteWorker(conn, queueName, deleteFunc)
}
