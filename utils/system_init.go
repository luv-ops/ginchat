package utils

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	Rdb *redis.Client
	Ctx = context.Background() //redis需要它
)

func InitConfig() {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	//为了使用docker使用了docker-compuse.yml的环境变量配置
	dsn := os.Getenv("MYSQL_DSN")
	redisAddr := os.Getenv("REDIS_ADDR")
	if dsn != "" {
		viper.Set("mysql.dsn", dsn)
	}
	if redisAddr != "" {
		viper.Set("redis.addr", redisAddr)
	}
	fmt.Println("config init success")
}
func InitMySql() {

	db, err := gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("mysql init success")
	DB = db
}

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
		fmt.Println("redis init fail")
		panic(err)
	}
	Rdb = rdb
	fmt.Println("redis init success")
}
