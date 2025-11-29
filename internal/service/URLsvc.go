package service

import (
	"biturl/internal/domain"
	"biturl/internal/dto"
	"biturl/internal/helper"
	"biturl/internal/queue/rabbitmq"
	queue "biturl/internal/queue/rabbitmq/deletetask"
	"biturl/internal/queue/rabbitmq/insertstat"
	"biturl/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type URLsvc struct {
	PG           repository.PGrepo
	RDB          repository.RedisRepo
	ClkhouseConn repository.ClkHouseRepo
	RabbitConn   *amqp.Connection
}

func (r URLsvc) CreateShortURL(input dto.URLdto, ctx context.Context) (string, string, error) {

	// generating expiry date for postgres db: 30 days
	pgTimeExpiresIn := helper.GenerateDate(30)

	// generating expiry date for redis: 3 days
	rdsTimeExpiresIn := helper.GenerateDate(3)

	// converting generated to time.until.
	ttl := time.Until(*rdsTimeExpiresIn)

	if ttl < 0 {
		ttl = 0
	}

	var shortCode string

	// checking if user provides original url
	if len(input.OriginalURL) == 0 {
		return "", "", errors.New("provide a url")
	}

	// if user provides short code, set it as short code
	if input.ShortCode != "" {
		shortCode = input.ShortCode
		// check if short code provided by user exists, throw error if it does...
		strcode, _ := r.PG.LoadURL(shortCode)
		fmt.Printf("strCode:%v", strcode)
		if strcode.ID != 0 {
			return "", "", errors.New("short url is not available")
			// else check if the length is less than 16 before saving it to db
		} else {
			fmt.Printf("code length: %v", len(input.ShortCode))
			if len(input.ShortCode) > domain.MaxLength {
				return "", "", errors.New("provide a shorter url or let us do it for you.")
			}
		}
	} else {
		// if user provides no short code, generate short code with snowflake id
		shortCode = helper.GenerateShortCode()
	}

	fmt.Printf("shortcode: %v", shortCode)

	statsAccessKey := "st_" + helper.GenerateShortCode()

	err := r.PG.CreateShortURL(&domain.URL{
		OriginalURL:    input.OriginalURL,
		ExpiresAt:      pgTimeExpiresIn,
		ShortCode:      shortCode,
		StatsAccessKey: statsAccessKey,
	})

	if err != nil {
		fmt.Printf("url creation error: %v", err)
		return "", "", errors.New("invalid request")
	}

	// save the short code and original url into redis for fast access. set time to live to 3 days
	// mutat
	input.ShortCode = shortCode
	err = r.RDB.CacheURL(ctx, input, ttl)

	if err != nil {
		fmt.Printf("saving code to redis error: %v", err)
	}

	return shortCode, statsAccessKey, nil
}

func (r URLsvc) LoadURL(shortCode string, ctx context.Context, stats repository.Stats) (string, error) {

	if len(shortCode) == 0 {
		return "", errors.New("no short code found")
	}

	val, err := r.RDB.GetCachedURL(ctx, shortCode)

	// checking if url does not exist in redis
	if err == redis.Nil {
		// then load if from postgres
		url, err := r.PG.LoadURL(shortCode)
		if err != nil {
			return "", errors.New("url not found")
		} else {
			// if url exists in postgres, create input dto struct since redis "cache url" function requires type of dto.URLdto. Then modify the struct so as to save it in redis. also create a time to live(ttl) for expiry date.
			var input dto.URLdto
			ttl := time.Until(*helper.GenerateDate(3))

			if ttl < 0 {
				ttl = 0
			}

			input.OriginalURL = url.OriginalURL
			input.ShortCode = shortCode
			r.RDB.CacheURL(ctx, input, ttl)

			// inserting stats to clickhouse db after redirect
			err := insertstat.PublishInsertStat(r.RabbitConn, rabbitmq.InsertClickhouseQueueKey, stats)

			if err != nil {
				fmt.Println("warning: could not enqueue clickhouse insert task", err)
			}

			return url.OriginalURL, nil
		}
	} else if err != nil {
		fmt.Printf("redis get error: %v", err)
		return "", errors.New("server error")
	} else {
		err := insertstat.PublishInsertStat(r.RabbitConn, rabbitmq.InsertClickhouseQueueKey, stats)
		if err != nil {
			fmt.Println("could not insert stats")
		}

		return val, nil
	}
}

func (r URLsvc) DeleteURL(shortcode string, ctx context.Context) error {

	if len(shortcode) == 0 {
		return errors.New("short code not found")
	}

	_, err := r.PG.LoadURL(shortcode)
	if err != nil {
		fmt.Println("cannot find url:", err)
		return errors.New("url not found")
	}

	err = r.PG.DeleteURL(shortcode)
	if err != nil {
		fmt.Println("error deleting url from postgres:", err)
		return errors.New("error deleting url")
	}

	if err := queue.PublishDeleteTask(r.RabbitConn, rabbitmq.DeleteRedisQueueKey, shortcode); err != nil {
		fmt.Println("warning: could not enqueue redis delete task", err)
	}

	return nil
}
