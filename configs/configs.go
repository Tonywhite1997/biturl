package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT         string
	DSN          string
	REDIS_ADDR   string
	RABBITMQ_URL string
}

func StartEnv() (cfg *Config, err error) {
	err = godotenv.Load()

	port := os.Getenv("PORT")
	dsn := os.Getenv("DSN")
	redisAddr := os.Getenv("REDIS_ADDR")
	rbtmq_url := os.Getenv("RABBITMQ_URL")

	return &Config{PORT: port, DSN: dsn, REDIS_ADDR: redisAddr, RABBITMQ_URL: rbtmq_url}, err
}
