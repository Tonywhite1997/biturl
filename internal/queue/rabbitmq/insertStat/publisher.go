package insertstat

import (
	"biturl/internal/repository"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishInsertStat(conn *amqp.Connection, queuename string, stats repository.Stats) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, _ := json.Marshal(repository.Stats{
		Id:           stats.Id,
		Url_short_id: stats.Url_short_id,
		User_ip:      stats.User_ip,
		User_agent:   stats.User_agent,
		Referer:      stats.Referer,
		Country:      stats.Country,
		City:         stats.City,
		Device:       stats.Device,
		OS:           stats.OS,
		Browser:      stats.Browser,
		Timestamp:    stats.Timestamp,
	})

	return ch.Publish(
		"",
		queuename,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}
