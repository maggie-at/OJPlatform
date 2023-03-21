package test

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

const (
	RedisAddr = "localhost:6379"
	//RedisPwd  = "alantam."
	RedisDb = 0
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr: RedisAddr,
	//Password: RedisPwd,
	DB: RedisDb,
})

func TestRedisSet(t *testing.T) {
	rdb.Set(ctx, "name", "alan", time.Second*60)
}

func TestRedisGet(t *testing.T) {
	result, err := rdb.Get(ctx, "21S151155@stu.hit.edu.cn").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}
