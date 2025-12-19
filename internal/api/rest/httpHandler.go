package rest

import (
	"biturl/internal/helper/geo"
	ratelimiter "biturl/internal/middleware/rate-limiter"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RestHandler struct {
	App            *fiber.App
	DB             *gorm.DB
	RDB            *redis.Client
	RabbitConn     *amqp.Connection
	ClickhouseConn clickhouse.Conn
	GEODB          *geo.GeoRedisCache
	StatsRatelimit *ratelimiter.RateLimiter
	URLRateLimit   *ratelimiter.RateLimiter
}
