package api

import (
	"biturl/configs"
	"biturl/internal/api/rest"
	"biturl/internal/api/rest/handlers"
	"biturl/internal/domain"
	"biturl/internal/helper"
	"biturl/internal/helper/geo"
	ratelimiter "biturl/internal/middleware/rate-limiter"
	"biturl/internal/queue/rabbitmq"
	"biturl/internal/repository"
	"biturl/internal/worker"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	c := cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Authorization , Content-Type, Accept",
		AllowMethods:     "PUT, PATCH,POST, GET, OPTIONS,DELETE",
		AllowCredentials: true,
	})

	app := fiber.New()

	app.Use(c)

	// postgres configuration
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("cannot connect to database")
	}
	// redis configuration
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.REDIS_ADDR,
	})

	geodb := geo.InitGeoDB("assets/geolocation/GeoLite2-City.mmdb", rdb)

	// fmt.Println(geodb.GEODB.City(net.ParseIP("8.8.8.8")))

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
	urlRateLimit := ratelimiter.NewRateLimiter(rdb, 10, time.Minute*1)
	statsRateLimit := ratelimiter.NewRateLimiter(rdb, 100, time.Minute*1)

	redisRepo := repository.NewRedisRepo(rdb)

	clkhouserepo := repository.NewClkHouseRepo(clkhouse)

	conn, err := amqp.Dial(cfg.RABBITMQ_URL)
	helper.FailOnError(err, "failed to connect to rabbitmq")

	// Delete worker
	go func() {
		defer helper.RecoverWorker()
		worker.StartDeleteWorker(conn, rabbitmq.DeleteRedisQueueKey, redisRepo.DeleteCachedURL)
	}()

	go func() {
		defer helper.RecoverWorker()
		worker.StartDeleteWorker(conn, rabbitmq.DeleteClickhouseStatQueueKey, clkhouserepo.DeleteStatsRecord)
	}()

	// Insert worker
	go func() {
		defer helper.RecoverWorker()
		worker.StartInsertWorker(conn, rabbitmq.InsertClickhouseQueueKey, clkhouserepo.Insert)
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
		ClickhouseConn: clkhouse,
		GEODB:          geodb,
		StatsRatelimit: statsRateLimit,
		URLRateLimit:   urlRateLimit,
	}
	handlers.SetupURLroutes(rh)
	handlers.SetupStatsRoute(rh)

	app.Listen(cfg.PORT)
	fmt.Printf("server running on port: %v", cfg.PORT)
}
