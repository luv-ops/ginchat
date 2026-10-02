package models

import "time"

type Friends struct {
	ID       uint      `json:"id"`
	UserId   uint      `json:"userId"`
	FriendId uint      `json:"friendId"`
	Status   uint      `json:"status"` // 1 正常 2 删除
	CreateAt time.Time `json:"createAt" gorm:"precision:0;autoCreateTime"`
}

// BlackList 黑名单功能  ，后期加
type BlackList struct {
	ID       uint      `json:"id"`
	UserId   uint      `json:"userId"`
	TargetId uint      `json:"targetId"`
	CreateAt time.Time `json:"createAt" gorm:"precision:0;autoCreateTime"`
}
type FriendResp struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`   // 申请人昵称
	Avatar   string `json:"avatar"` // 申请人头像
	IsOnline int    `json:"isOnline"`
}
type FriendReq struct {
	ID       uint      `json:"id"`
	FromId   uint      `json:"fromId"`
	TargetId uint      `json:"targetId"`
	Status   uint      `json:"status"` //0:申请中 1:已同意 2:已拒绝 3:已读不处理
	CreateAt time.Time `json:"createAt" gorm:"precision:0;autoCreateTime"`
}
type FriendApplyResp struct {
	FromId   uint   `json:"fromId"`
	TargetId uint   `json:"targetId"`
	Type     string `json:"type"`
	Name     string `json:"name"`   // 申请人昵称
	Avatar   string `json:"avatar"` // 申请人头像
	Msg      string `json:"msg"`    // 申请附言
	Status   int    `json:"status"`
	CreateAt string `json:"create_at"`
}
