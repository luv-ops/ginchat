package MQ

import (
	"fmt"

	"github.com/goccy/go-json"
)

const (
	TopicPrivateMsg  = "im_private_topic" // 私聊消息topic
	TopicGroupMsg    = "im_group_topic"   // 群聊消息topic
	TopicFriendReq   = "im_friend_req_topic"
	TopicOfflinePush = "im_offline_push_topic" // 离线推送topic
	ConsumerGroupID  = "im-consumer-group"     // 消费组ID
	DlqTopic         = "im_dlq_topic"          // 死信队列
	MaxRetryCount    = 3                       // 最大重试次数

	// 业务聊天类型常量
	ChatTypePrivate       = 1
	ChatTypeGroup         = 2
	ChatTypeFriendRequest = 3
)

// kafka消息数据结构
type MsgDTO struct {
	MsgID      string `json:"msg_id"`
	FromID     uint   `json:"from_uid"`
	TargetID   uint   `json:"target_id"`
	ChatType   int    `json:"chat_type"`
	MsgType    int    `json:"msg_type"`
	Content    string `json:"content"`
	SendTime   int64  `json:"send_time"`
	UserOnline bool   `json:"user_online"`
}

// 根据消息的type来分区，保证消息有序
func (m *MsgDTO) GetPartitionKey() string {
	if m.ChatType == ChatTypePrivate {
		return fmt.Sprintf("%d_%d", min(m.FromID, m.TargetID), max(m.FromID, m.TargetID))
	} else if m.ChatType == ChatTypeGroup {
		return fmt.Sprintf("g_%d", m.TargetID)
	} else {
		return fmt.Sprintf("f_%d", m.FromID)
	}

}
func (m *MsgDTO) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
func (m *MsgDTO) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}

// MessageHandler 由chatService实现 需要注入chatService
type MessageHandler interface {
	HandleMsg(dto *MsgDTO) error
}

var MsgHandler MessageHandler

// FriendHandler 由friendService实现
type FriendHandler interface {
	HandleFReq(dto *MsgDTO) error
}

var FriReqHandler FriendHandler
