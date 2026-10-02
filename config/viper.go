package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
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
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if dsn != "" {
		viper.Set("Mysql.dsn", dsn)
	}
	if redisAddr != "" {
		viper.Set("redis.addr", redisAddr)
	}
	if kafkaBroker != "" {
		viper.Set("kafka.broker", kafkaBroker)
	}
	log.Println("config Autowired success")
}
