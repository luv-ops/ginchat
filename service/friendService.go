package service

import (
	"GinChat/MQ"
	"GinChat/mapper"
	"GinChat/models"
	"GinChat/redis"
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FriendService struct {
	friendMapper  *mapper.FriendMapper
	userMapper    *mapper.UserMapper
	messageSender IMessageSender
	db            *gorm.DB
	kafkaCli      *MQ.KafkaClient
}

func NewFriendService(fm *mapper.FriendMapper, uM *mapper.UserMapper, mS IMessageSender,
	db *gorm.DB, kC *MQ.KafkaClient) *FriendService {
	return &FriendService{
		friendMapper:  fm,
		userMapper:    uM,
		messageSender: mS,
		db:            db,
		kafkaCli:      kC,
	}
}

// 供controller层使用
func (s *FriendService) AddFriend(ctx context.Context, friendReq *models.FriendReq) error {
	//组装kafka消息
	dto := MQ.MsgDTO{
		MsgID:    uuid.NewString(),
		FromID:   friendReq.FromId,
		TargetID: friendReq.TargetId,
		ChatType: MQ.ChatTypeFriendRequest,
	}
	return s.kafkaCli.SendFriendReq(ctx, &dto)

}

// 供kafka异步调用，并通过实现防止循环引用
func (s *FriendService) HandleFReq(dto *MQ.MsgDTO) error {
	//	组装friendReq
	friendReq := &models.FriendReq{
		FromId:   dto.FromID,
		TargetId: dto.TargetID,
	}
	var count int64
	//查询用户是否存在
	err := s.userMapper.UserExistById(friendReq.TargetId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	err = s.friendMapper.FriendReqExist(friendReq.FromId, friendReq.TargetId, &count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("好友请求已经存在,请等待对方回应")
	}
	//判断是否已经是好友
	var friendCount int64
	err = s.friendMapper.FriendsExist(friendReq.FromId, friendReq.TargetId, &friendCount)
	if err != nil {
		return err
	}
	if friendCount > 0 {
		return errors.New("你们已经是好友")
	}
	err = s.friendMapper.CreateFriendReq(friendReq)
	if err != nil {
		return err
	}
	//推送添加好友申请
	message := models.Message{
		FromId:   friendReq.FromId,
		TargetId: friendReq.TargetId,
		Type:     "friendRequest",
	}
	go func() {
		if e := recover(); e != nil {
			fmt.Println("自增好友请求未读缓存更新pinic", e)
		}
		err = redis.IncrFriendReqUnread(friendReq.TargetId)
		if err != nil {
			fmt.Println("自增好友请求未读缓存失败", err.Error())
		}
	}()
	return s.messageSender.SendWs(&message)

}

//func (s *FriendService) AddFriend(friendReq *models.FriendReq) error {
//	var count int64
//	//查询用户是否存在
//	err := s.userMapper.UserExistById(friendReq.TargetId)
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return errors.New("用户不存在")
//		}
//		return err
//	}
//
//	err = s.friendMapper.FriendReqExist(friendReq.FromId, friendReq.TargetId, &count)
//	if err != nil {
//		return err
//	}
//	if count > 0 {
//		return errors.New("好友请求已经存在,请等待对方回应")
//	}
//	//判断是否已经是好友
//	var friendCount int64
//	err = s.friendMapper.FriendsExist(friendReq.FromId, friendReq.TargetId, &friendCount)
//	if err != nil {
//		return err
//	}
//	if friendCount > 0 {
//		return errors.New("你们已经是好友")
//	}
//	err = s.friendMapper.CreateFriendReq(friendReq)
//	if err != nil {
//		return err
//	}
//	//推送添加好友申请
//	message := models.Message{
//		FromId:   friendReq.FromId,
//		TargetId: friendReq.TargetId,
//		Type:     "friendRequest",
//	}
//	go func() {
//		if e := recover(); e != nil {
//			fmt.Println("自增好友请求未读缓存更新pinic", e)
//		}
//		err = redis.IncrFriendReqUnread(friendReq.TargetId)
//		if err != nil {
//			fmt.Println("自增好友请求未读缓存失败", err.Error())
//		}
//	}()
//	return s.messageSender.SendWs(&message)
//
//}

func (s *FriendService) RequestList(targetId uint) ([]models.FriendApplyResp, error) {
	list := []models.FriendApplyResp{}
	//查全部，
	err := s.friendMapper.SelectFriendReqListAndInfo(targetId, &list)
	if err != nil {
		return list, err
	}
	return list, nil
}

func (s *FriendService) Accept(fromId uint, targetId uint) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := s.friendMapper.UpdateStatusWithTx(tx, fromId, targetId)
		if err != nil {
			return err
		}
		err = s.friendMapper.CreateFriendsWithTx(tx, fromId, targetId)
		if err != nil {
			return err
		}
		err = s.friendMapper.CreateFriendsWithTx(tx, targetId, fromId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	//同意好友申请后，双向缓存失效
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("panic:", r)
			}
		}()
		key1 := redis.KeyFriendList + strconv.Itoa(int(fromId))
		key2 := redis.KeyFriendList + strconv.Itoa(int(targetId))
		_, err = redis.Rdb.Del(redis.Ctx, key1).Result()
		if err != nil {
			fmt.Println("redis error:", err.Error())
		}
		_, err = redis.Rdb.Del(redis.Ctx, key2).Result()
		if err != nil {
			fmt.Println("redis error:", err.Error())
		}

	}()
	return nil
}
func (s *FriendService) Reject(fromId uint, targetId uint) error {
	return s.friendMapper.UpdateStatus(fromId, targetId)
}

func (s *FriendService) GetFriendList(id uint) ([]models.FriendResp, error) {
	list := []models.FriendResp{}
	//先查redis
	list, err := redis.GetFriendList(id)
	if err == nil && len(list) > 0 {
		//更新状态，此字段不存redis，单独维护
		for i, friend := range list {
			status, _ := redis.GetUserLine(friend.Id)
			list[i].IsOnline = status
		}
		return list, nil
	}
	err = s.friendMapper.SelectFriendListAndInfo(id, &list)

	if err != nil {
		return list, err
	}
	go func() {
		if e := recover(); e != nil {
			fmt.Println("pinic", e)
		}
		err = redis.SaveFriendList(id, list)
		if err != nil {
			fmt.Println("更新好友列表缓存失败", err.Error())
		}
	}()

	//更新状态，此字段不存redis，单独维护
	for i, friend := range list {
		status, _ := redis.GetUserLine(friend.Id)
		list[i].IsOnline = status
	}
	return list, nil
}

func (s *FriendService) UnReadCount(userid uint) (int64, error) {
	var count int64
	//先查redis
	unread, err := redis.GetFriendReqUnread(userid)
	if err == nil {
		return unread, nil
	}
	err = s.friendMapper.FriendReqUnreadCount(userid, &count)
	//写回redis
	go func() {
		if e := recover(); e != nil {
			fmt.Println("好友请求未读数量接口pinic", e)
		}
		err = redis.SetFriendReqUnread(userid, count)
		if err != nil {
			fmt.Println("更新好友请求未读数量失败", err.Error())
		}
	}()
	if err != nil {
		return count, err
	}
	return count, nil

}

func (s *FriendService) HasRead(userId uint) error {
	err := s.friendMapper.FriendReqHasRead(userId)
	if err != nil {
		return err
	}
	_ = redis.SetFriendReqUnread(userId, 0)
	return nil

}
