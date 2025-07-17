package redis

import (
	"context"
	"desafio_rate-limiter/config"
	"desafio_rate-limiter/internal/storage"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore(cfg config.Config) (storage.RateLimiterStore, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("erro ao conectar ao Redis: %w", err)
	}

	return &RedisStore{
		client: rdb,
		ctx:    ctx,
	}, nil
}

func (r *RedisStore) Increment(key string) (int, error) {
	count, err := r.client.Incr(r.ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *RedisStore) SetExpiration(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

func (r *RedisStore) Get(key string) (int, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *RedisStore) Reset(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
