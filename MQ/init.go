package MQ

import (
	"context"
	"log"

	"github.com/spf13/viper"
)

var GlobalKafkaCli *KafkaClient

func InitKafkaConfig() {
	broker := viper.GetString("kafka.broker")
	client, err := NewKafkaClient([]string{broker})
	if err != nil {
		log.Println("kafka初始化失败", err.Error())
		panic(err)
	}
	GlobalKafkaCli = client
	log.Println("kafka初始化成功")
}

// StartAllConsumers 统一启动所有消费者协程
func StartAllConsumers(ctx context.Context) {
	// websocket相关消息消费
	go GlobalKafkaCli.StartCommonConsumer(ctx, TopicPrivateMsg)
	go GlobalKafkaCli.StartCommonConsumer(ctx, TopicGroupMsg)
	go GlobalKafkaCli.StartCommonConsumer(ctx, TopicFriendReq)

	log.Println("✅ 全部Kafka消费者后台启动完成")
}

// CloseKafka 优雅释放kafka资源，main收到停机信号调用
func CloseKafka() {
	if GlobalKafkaCli == nil {
		return
	}
	if err := GlobalKafkaCli.Close(); err != nil {
		log.Printf("关闭Kafka连接异常: %v", err)
	}
	log.Println("✅ Kafka资源全部释放")
}
