package models

import (
	"github.com/go-redis/redis/v8"
)

const (
	RedisAddr = "localhost:6379"
	RedisDb   = 0
)

var RDB = InitRedis()

func InitRedis() *redis.Client {
	var rdbCli = redis.NewClient(&redis.Options{
		Addr: RedisAddr,
		DB:   RedisDb,
		//Password: RedisPwd,
	})
	return rdbCli
}
