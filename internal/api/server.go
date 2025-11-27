package api

import (
	"biturl/configs"
	"biturl/internal/api/rest"
	"biturl/internal/api/rest/handlers"
	"biturl/internal/domain"
	"biturl/internal/helper"
	ratelimiter "biturl/internal/middleware/rate-limiter"
	"biturl/internal/queue/rabbitmq"
	"biturl/internal/repository"
	"biturl/internal/worker"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func StartsServer() {

	cfg, err := configs.StartEnv()

	ctx := context.Background()

	if err != nil {
		log.Fatal("error loading your environment variables")
	}

	app := fiber.New()

	// postgres configuration
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("cannot connect to database")
	}
	// redis configuration
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.REDIS_ADDR,
	})

	// clickhouse db configuration
	clkhouse, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.CLICKHOUSE_ADDR},

		Auth: clickhouse.Auth{
			Database: cfg.CLICKHOUSE_DB,
			Username: cfg.CLICKHOUSE_USER,
			Password: cfg.CLICKHOUSE_PASSWORD,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
	})

	if err != nil {
		log.Fatal("clickhouse connection error: ", err)
	}

	if err := clkhouse.Ping(ctx); err != nil {
		log.Fatal("clickhouse ping failed: ", err)
	}

	// rate limiting configuration
	rl := ratelimiter.NewRateLimiter(rdb, 10, time.Minute)
	app.Use(rl.Middleware())

	redisRepo := repository.NewRedisRepo(rdb)

	conn, err := amqp.Dial(cfg.RABBITMQ_URL)
	helper.FailOnError(err, "failed to connect to rabbitmq")

	ch, err := conn.Channel()
	helper.FailOnError(err, "failed to open a channel")

	_, err = ch.QueueDeclare(rabbitmq.DeleteRedisQueueKey, true, false, false, false, nil)
	helper.FailOnError(err, "faild to declare queue")

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("worker panicked", r)
			}
		}()
		worker.StartDeleteWorker(ch, rabbitmq.DeleteRedisQueueKey, redisRepo.DeleteCachedURL)
	}()

	err = db.AutoMigrate(&domain.URL{})

	if err != nil {
		fmt.Println(err)
		log.Fatal("error migrating your db model")
	}

	if err := helper.InitializeSnowflake(1); err != nil {
		log.Fatal(err)
	}

	rh := &rest.RestHandler{
		App:            app,
		DB:             db,
		RDB:            rdb,
		RabbitConn:     conn,
		ClickhouseConn: &clkhouse,
	}
	handlers.SetupURLroutes(rh)

	app.Listen(cfg.PORT)
	fmt.Printf("server running on port: %v", cfg.PORT)
}
