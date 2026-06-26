package MQ

import (
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaClient 全局统一操作对象
type KafkaClient struct {
	brokers []string      // kafka集群地址
	writer  *kafka.Writer // 正常消息生产者
	dlqW    *kafka.Writer // 死信队列生产者
}

func NewKafkaClient(brokers []string) (*KafkaClient, error) {

	//主消息写入器
	mainWriter := kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.Hash{},
		WriteTimeout: 300 * time.Millisecond,
		RequiredAcks: kafka.RequireOne, // 首领分区写入成功就返回,对于IM项目非常适合，可靠性和快速
		BatchSize:    16 * 1024,        //开启缓存池，提升处理速度
		BatchTimeout: 5 * time.Millisecond,
	}
	//死信队列写入器
	dlqWriter := kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.Hash{},
		Topic:        DlqTopic,
		WriteTimeout: 300 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
		BatchSize:    16 * 1024, //开启缓存池，提升处理速度
		BatchTimeout: 5 * time.Millisecond,
	}
	return &KafkaClient{
		brokers: brokers,
		writer:  &mainWriter,
		dlqW:    &dlqWriter,
	}, nil
}

func (k *KafkaClient) Close() error {
	_ = k.dlqW.Close()
	return k.writer.Close()
}
func (k *KafkaClient) PreCreateTopic() error {
	conn, err := kafka.Dial("tcp", k.brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()
	topics := []kafka.TopicConfig{
		{Topic: TopicPrivateMsg, NumPartitions: 1, ReplicationFactor: 1},
		{Topic: TopicGroupMsg, NumPartitions: 1, ReplicationFactor: 1},
		{Topic: TopicFriendReq, NumPartitions: 1, ReplicationFactor: 1},
		{Topic: DlqTopic, NumPartitions: 1, ReplicationFactor: 1}, // 别忘了死信队列
	}
	err = conn.CreateTopics(topics...)
	if err != nil {
		if errors.Is(err, kafka.TopicAlreadyExists) {
			return nil
		}
		return err
	}
	return nil
}
