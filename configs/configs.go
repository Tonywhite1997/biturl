package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT                string
	DSN                 string
	REDIS_ADDR          string
	RABBITMQ_URL        string
	CLICKHOUSE_ADDR     string
	CLICKHOUSE_DB       string
	CLICKHOUSE_PASSWORD string
	CLICKHOUSE_USER     string
}

func StartEnv() (cfg *Config, err error) {
	err = godotenv.Load()

	port := os.Getenv("PORT")
	dsn := os.Getenv("DSN")
	redisAddr := os.Getenv("REDIS_ADDR")
	rbtmq_url := os.Getenv("RABBITMQ_URL")
	clickhouse_addr := os.Getenv("CLICKHOUSE_ADDR")
	clickhouse_db := os.Getenv("CLICKHOUSE_DB")
	clickhouse_user := os.Getenv("CLICKHOUSE_USER")
	clickhouse_password := os.Getenv("CLICKHOUSE_PASSWORD")

	return &Config{PORT: port, DSN: dsn, REDIS_ADDR: redisAddr, RABBITMQ_URL: rbtmq_url, CLICKHOUSE_ADDR: clickhouse_addr, CLICKHOUSE_DB: clickhouse_db, CLICKHOUSE_PASSWORD: clickhouse_password, CLICKHOUSE_USER: clickhouse_user}, err
}
