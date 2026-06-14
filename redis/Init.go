package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     "",
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.PoolSize"),
		MinIdleConns: viper.GetInt("redis.MinIdleConns"),
	})
	_, cancel := context.WithTimeout(Ctx, 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		fmt.Println("redis Autowired fail")
		panic(err)
	}
	Rdb = rdb
	log.Println("redis Autowired success")
}
