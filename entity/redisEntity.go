package entity

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var ctx = context.Background()

type RedisEntity interface {
	Set(key string, value interface{}, expire time.Duration) (string, error)
	Get(key string) (interface{}, error)
	Del(key string) (interface{}, error)
	GetInt(key string) (int, error)
	IncrBy(key string, value int64) (uint64, error)
	ExpireAt(key string, time time.Time) bool
}

type redisConnection struct {
	connection *redis.Client
}

func NewRedisEntity(rdb *redis.Client) RedisEntity {
	return &redisConnection{
		connection: rdb,
	}
}

func (rdb *redisConnection) Set(key string, value interface{}, expire time.Duration) (string, error) {
	val, err := rdb.connection.Set(ctx, key, value, expire).Result()
	return val, err
}

func (rdb *redisConnection) Get(key string) (interface{}, error) {
	val, err := rdb.connection.Get(ctx, key).Result()
	return val, err
}

func (rdb *redisConnection) Del(key string) (interface{}, error) {
	val, err := rdb.connection.Del(ctx, key).Result()
	return val, err
}

func (rdb *redisConnection) GetInt(key string) (int, error) {
	val, err := rdb.connection.Get(ctx, key).Int()
	return val, err
}

// Increase the number of requests
func (rdb *redisConnection) IncrBy(key string, value int64) (uint64, error) {
	val, err := rdb.connection.IncrBy(ctx, key, value).Uint64()
	return val, err
}

// Expire time
func (rdb *redisConnection) ExpireAt(key string, time time.Time) bool {
	val := rdb.connection.ExpireAt(ctx, key, time).Val()
	return val
}
