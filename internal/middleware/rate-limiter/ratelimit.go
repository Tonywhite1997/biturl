package ratelimiter

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	RedisClient *redis.Client
	Limit       int
	Period      time.Duration
}

func NewRateLimiter(rdb *redis.Client, limit int, period time.Duration) *RateLimiter {
	return &RateLimiter{
		RedisClient: rdb,
		Limit:       limit,
		Period:      period,
	}
}

func (rl *RateLimiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := fmt.Sprintf("rate:%s", c.IP())
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		// Increment the counter
		count, err := rl.RedisClient.Incr(ctx, key).Result()
		if err != nil {
			fmt.Printf("Redis error: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "server error",
			})
		}

		// If this is the first request, set the TTL for the period
		if count == 1 {
			_, _ = rl.RedisClient.Expire(ctx, key, rl.Period).Result()
		}

		// Compute remaining requests
		remaining := rl.Limit - int(count)
		if remaining < 0 {
			remaining = 0
		}

		// Get remaining TTL for Retry-After header
		ttl, err := rl.RedisClient.TTL(ctx, key).Result()
		if err != nil || ttl < 0 {
			ttl = rl.Period
		}

		// Set standard rate-limit headers
		c.Set("X-RateLimit-Limit", strconv.Itoa(rl.Limit))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))

		// Reject if limit exceeded
		if count > int64(rl.Limit) {
			c.Set("Retry-After", strconv.FormatInt(int64(ttl.Seconds()), 10))
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "rate limit exceeded",
			})
		}

		return c.Next()
	}
}
