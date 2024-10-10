package redis

import (
	"ebank/shared/microservice"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type Redis interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
}

func New() Redis {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", microservice.RedisHost, microservice.RedisPort),
	})
	return client
}
