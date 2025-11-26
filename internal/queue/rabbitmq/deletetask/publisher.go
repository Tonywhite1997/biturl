package queue

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type DeleteTask struct {
	ShortCode string `json:"short_code"`
}

func PublishDeleteTask(conn *amqp.Connection, queuename, shortCode string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, _ := json.Marshal(DeleteTask{ShortCode: shortCode})

	return ch.Publish(
		"",
		queuename,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
