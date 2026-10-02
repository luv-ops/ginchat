package MQ

import (
	"fmt"

	"github.com/goccy/go-json"
)

const (
	ConsumerGroupID = "im-consumer-group" // 消费组ID
	MaxRetryCount   = 3                   // 最大重试次数
	// 业务聊天类型常量
	ChatTypePrivate = iota
	ChatTypeGroup
)

// Topic
const (
	TopicPrivateMsg = "im_private_topic" // 私聊消息topic
	TopicGroupMsg   = "im_group_topic"   // 群聊消息topic
	DlqTopic        = "im_dlq_topic"     // 死信队列
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
type GroupDTO struct {
	GroupID   uint   `json:"group_id"`
	GroupName string `json:"group_name"`
	OwnerID   uint   `json:"owner_id"`
	InviteIds []uint `json:"invite_ids"`
	Type      int    `json:"type"`
}
type UserDTO struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// GetPartitionKey 根据消息的type来分区，保证消息有序
func (m *MsgDTO) GetPartitionKey() string {
	if m.ChatType == ChatTypePrivate {
		return fmt.Sprintf("%d_%d", min(m.FromID, m.TargetID), max(m.FromID, m.TargetID))
	} else if m.ChatType == ChatTypeGroup {
		return fmt.Sprintf("g_%d", m.TargetID)
	}
	return ""

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
