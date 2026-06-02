package models

import (
	"time"
)

type Message struct {
	ID       uint      `json:"id"`
	FromId   uint      `json:"fromId"`
	TargetId uint      `json:"targetId"`
	Type     string    `json:"type"`  //单聊，群聊 chat ,groupMessage,friendRequest
	Media    int       `json:"media"` //图片，文本，音频
	Content  string    `json:"content"`
	MsgType  int       `json:"msgType"` // 0 文本 1 图片  (后续可以扩展为2 音频)
	CreateAt time.Time `json:"createAt" gorm:"precision:0;autoCreateTime"`
}
type MessageVO struct {
	Message
	FromName   string `json:"name"`
	FromAvatar string `json:"avatar"`
}
type FriendReqMsg struct {
	FriendId uint      `json:"friendId"`
	TargetId uint      `json:"targetId"`
	Type     string    `json:"type"` //单聊，群聊 upload group,friendRequest
	CreateAt time.Time `json:"createAt" gorm:"precision:0;autoCreateTime"`
	Name     string    `json:"name"`
	Avatar   string    `json:"avatar"`
	Status   uint      `json:"status"`
}

type MessageReq struct {
	PeerId uint   `form:"peerId" json:"peerId"`
	Page   int    `form:"page" json:"page"`
	Size   int    `form:"size" json:"size"`
	Type   string `form:"type" json:"type"` //chat单聊 groupMessage群聊
}
