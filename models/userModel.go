package models

import "time"

type UserBasic struct {
	ID            uint
	Name          string `json:"name" gorm:"not null"`
	Password      string `gorm:"not null"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string
	Salt          string
	Avatar        string `json:"avatar"`
	ClientIp      string
	ClientPort    string
	LoginTime     time.Time `gorm:"precision:0;default:null"`
	HeartbeatTime time.Time `gorm:"precision:0;default:null"`
	LoginOutTime  time.Time `gorm:"precision:0;default:null"`
	IsLogout      bool
	DeviceInfo    string
	CreateAt      time.Time  `gorm:"precision:0;autoCreateTime"`
	UpdateAt      time.Time  `gorm:"precision:0;autoCreateTime"`
	DeleteAt      *time.Time `gorm:"precision:0;default:null"`
}

type RegisterReq struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}
type LoginReq struct {
	Name     string
	Password string
}

// 过滤用户信息中敏感字段
type UserRes struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Token  string `json:"token"`
	Avatar string `json:"avatar"`
}
type UpdateReq struct {
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type OnlineRes struct {
	Status int `json:"status"`
}

func (u UserBasic) TableName() string {
	return "user_basic"
}
