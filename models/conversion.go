package models

import "time"

type Conversation struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"userID"` // 当前用户
	PeerID      uint      `json:"peerID"` // 对方ID
	LastMsg     string    `json:"lastMsg"`
	LastTime    time.Time `json:"lastTime" gorm:"precision:0;autoCreateTime"`
	UnreadCount uint      `json:"unreadCount"`
	Type        int       `json:"type"` // 0单聊 1群聊
}

type ConversationInfo struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"userId"` // 当前用户
	PeerID      uint      `json:"peerId"` // 对方ID
	LastMsg     string    `json:"lastMsg"`
	LastTime    time.Time `json:"lastTime" gorm:"precision:0;autoCreateTime"`
	UnreadCount uint      `json:"unreadCount"`
	Type        int       `json:"type"` // 0单聊 1群聊
	Name        string    `json:"name"`
	Avatar      string    `json:"avatar"`
}
