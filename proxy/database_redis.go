package proxy

import (
	"time"

	"github.com/go-redis/redis"
)

type DatabaseRedis struct {
	Client *redis.Client
}

func (db DatabaseRedis) Get(key string) (string, error) {
	return db.Client.Get(key).Result()
}

func (db DatabaseRedis) Set(key string, value string, expiration time.Duration) error {
	return db.Client.Set(key, value, expiration).Err()
}

func (db DatabaseRedis) Delete(key string) error {
	return db.Client.Del(key).Err()
}

func (db DatabaseRedis) Keys(prefix string) ([]string, error) {
	return db.Client.Keys(prefix + "*").Result()
}

func (db DatabaseRedis) Health() error {
	return db.Client.Ping().Err()
}
