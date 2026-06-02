package models

import "time"

// Group 群聊信息表
type GroupModel struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	GroupName  string    `gorm:"size:64;not null" json:"group_name"` // 群名称
	Avatar     string    `gorm:"size:255" json:"avatar"`             // 群头像
	Notice     string    `gorm:"size:1000" json:"notice"`            // 群公告
	TotalCount uint      `json:"total_count"`                        // 群成员数量
	OwnerID    uint      `gorm:"not null" json:"owner_id"`           // 群主ID
	IsAllMute  int       `gorm:"default:0" json:"is_all_mute"`       // 0=不禁言 1=全员禁言
	CreatedAt  time.Time `json:"created_at" gorm:"precision:0;autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"precision:0"`
}

// GroupMember 群成员表（包含角色、禁言状态）
type GroupMember struct {
	ID      uint `gorm:"primaryKey" json:"id"`
	GroupID uint `gorm:"index;not null" json:"group_id"`
	UserID  uint `gorm:"index;not null" json:"user_id"` // 用户ID
	// 角色 0=普通成员 1=管理员 2=群主
	Role int `gorm:"default:0" json:"role"`
	// 是否禁言 0=否 1=是
	IsMute int `gorm:"default:0" json:"is_mute"`
	// 禁言到期时间（永久禁言可设为 2099年）
	MuteEndTime time.Time `gorm:"default:null" json:"mute_end_time"`
	CreatedAt   time.Time `json:"created_at"`
}
type CreateGroupReq struct {
	GroupName string `json:"group_name"`
}
type InviteReq struct {
	GroupId   uint   `json:"group_id"`
	InvitedId []uint `json:"invited_id"`
}

// 群成员VO（带用户信息）
type GroupMemberVO struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Role   int    `json:"role"` // 2群主 1管理员 0普通成员
}
type GroupMemberReq struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

// 群详情VO
type GroupDetailVO struct {
	GroupID    uint            `json:"group_id"`
	TotalCount uint            `json:"totalCount"`
	GroupName  string          `json:"group_name"`
	Avatar     string          `json:"avatar"`
	Notice     string          `json:"notice"`
	Members    []GroupMemberVO `json:"members"` // 成员列表
}
