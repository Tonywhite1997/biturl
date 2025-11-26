package rest

import (
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RestHandler struct {
	App        *fiber.App
	DB         *gorm.DB
	RDB        *redis.Client
	RabbitConn *amqp.Connection
}
