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
	ChatTypeFriendRequest
	ChatTypeFriendRequestAccept
	ChatTypeFriendRequestHasRead
	GroupCreate
	GroupInvite
)

// Topic
const (
	TopicPrivateMsg       = "im_private_topic"             // 私聊消息topic
	TopicGroupMsg         = "im_group_topic"               // 群聊消息topic
	TopicFriendReq        = "im_friend_req_topic"          // 好友请求topic
	DlqTopic              = "im_dlq_topic"                 // 死信队列
	TopicFriendReqAccept  = "im_friend_req_accept_topic"   // 好友请求接受topic
	TopicFriendReqHasRead = "im_friend_req_has_read_topic" // 好友请求已读topic
	TopicGroupCreate      = "im_group_create_topic"        // 群聊创建topic
	TopicGroupInvite      = "im_group_invite_topic"        //群聊邀请topic
	TopicUserCreate       = "im_user_topic"                // 用户创建topic
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

// 根据消息的type来分区，保证消息有序
func (m *MsgDTO) GetPartitionKey() string {
	if m.ChatType == ChatTypePrivate {
		return fmt.Sprintf("%d_%d", min(m.FromID, m.TargetID), max(m.FromID, m.TargetID))
	} else if m.ChatType == ChatTypeGroup {
		return fmt.Sprintf("g_%d", m.TargetID)
	} else if m.ChatType == ChatTypeFriendRequest { // 好友请求消息
		return fmt.Sprintf("fq_%d", m.FromID)
	} else if m.ChatType == ChatTypeFriendRequestAccept { // 好友请求接受消息
		return fmt.Sprintf("fqA_%d", m.FromID)
	} else if m.ChatType == ChatTypeFriendRequestHasRead { // 好友请求已读消息
		return fmt.Sprintf("fqH_%d", m.FromID)
	}
	return ""

}
func (m *MsgDTO) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
func (m *MsgDTO) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}
func (m *GroupDTO) GetPartitionKey(Type int) string {
	if Type == GroupCreate {
		return fmt.Sprintf("gc_%d", m.GroupID)
	} else if Type == GroupInvite {
		return fmt.Sprintf("gi_%d", m.GroupID)
	}
	return ""
}
func (m *GroupDTO) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
func (m *GroupDTO) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}
func (m *UserDTO) GetPartitionKey() string {
	return fmt.Sprintf("uc_%s", m.Name)
}
func (m *UserDTO) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
func (m *UserDTO) Unmarshal(data []byte) error {
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
	HandleFReqAccept(dto *MsgDTO) error
	HandleFReqHasRead(dto *MsgDTO) error
}

var FriReqHandler FriendHandler

type GroupHandler interface {
	HandleGroupCreate(dto *GroupDTO) error
	HandleGroupInvite(dto *GroupDTO) error
}

var GroHandler GroupHandler

type UserHandler interface {
	HandleUserCreate(dto *UserDTO) error
}

var UHandler UserHandler
