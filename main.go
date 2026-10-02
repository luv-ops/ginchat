package main

import (
	"GinChat/Autowired"
	"GinChat/MQ"
	"GinChat/Mysql"
	"GinChat/config"
	"GinChat/redis"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	Mysql.InitMySql()
	redis.InitRedis()
	MQ.InitKafkaConfig()
	// 创建根 context，用于控制消费者生命周期
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	MQ.StartAllConsumers(ctx)
	//依赖注入
	route := Autowired.InitAll()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	route.Setup(r)
	r.Static("/static", "./static")
	// 优雅退出处理
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("收到退出信号，正在关闭...")
		cancel()        // 取消 context，停止所有消费者
		MQ.CloseKafka() // 释放 Kafka 资源
		os.Exit(0)
	}()
	r.Run(":8080")
}
