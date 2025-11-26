package repository

import (
	"biturl/internal/dto"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepo interface {
	CacheURL(ctx context.Context, input dto.URLdto, ttl time.Duration) error
	GetCachedURL(ctx context.Context, shortCode string) (string, error)
	DeleteCachedURL(ctx context.Context, shortCode string) error
}

type redisrepo struct {
	RDB *redis.Client
}

// CacheURL implements [RedisRepo].
func (r *redisrepo) CacheURL(ctx context.Context, input dto.URLdto, ttl time.Duration) error {
	return r.RDB.Set(ctx, input.ShortCode, input.OriginalURL, ttl).Err()
}

// GetCachedURL implements [RedisRepo].
func (r *redisrepo) GetCachedURL(ctx context.Context, shortCode string) (string, error) {
	val, err := r.RDB.Get(ctx, shortCode).Result()
	return val, err
}

// DeleteCachedURL implements [RedisRepo].
func (r *redisrepo) DeleteCachedURL(ctx context.Context, shortCode string) error {
	return r.RDB.Del(ctx, shortCode).Err()
}

func NewRedisRepo(rdb *redis.Client) RedisRepo {
	return &redisrepo{RDB: rdb}
}
