package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

type redisCache struct {
	rdb *redis.Client
}

func NewRedisCache(rdb *redis.Client) Cache {
	return &redisCache{rdb: rdb}
}

func (r *redisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(val, dest)
}

func (r *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, key, val, expiration).Err()
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}
