package MQ

import "github.com/spf13/viper"

const (
	TopicPrivateMsg  = "im_private_topic"      // 私聊消息
	TopicGroupMsg    = "im_group_topic"        // 群聊消息
	TopicOfflinePush = "im_offline_push_topic" // 离线推送
	ConsumerGroupID  = "im-consumer-group"     // 消费组
	DlqTopic         = "im_dlq_topic"          // 死信队列
	MaxRetryCount    = 3                       // 最大重试次数
)

var BrokersAddr = viper.GetString("kafka.broker")
