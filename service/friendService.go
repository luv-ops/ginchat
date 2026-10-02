package service

import (
	"GinChat/mapper"
	"GinChat/models"
	"GinChat/redis"
	"context"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type FriendService struct {
	friendMapper  *mapper.FriendMapper
	userMapper    *mapper.UserMapper
	messageSender IMessageSender
	db            *gorm.DB
}

func NewFriendService(fm *mapper.FriendMapper, uM *mapper.UserMapper, mS IMessageSender,
	db *gorm.DB) *FriendService {
	return &FriendService{
		friendMapper:  fm,
		userMapper:    uM,
		messageSender: mS,
		db:            db,
	}
}

// AddFriend 供controller层使用
func (s *FriendService) AddFriend(ctx context.Context, friendReq *models.FriendReq) error {
	var exist bool
	//查询用户是否存在
	err := s.userMapper.UserExistById(ctx, friendReq.TargetId, &exist)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("用户不存在")
	}

	err = s.friendMapper.FriendReqExist(ctx, friendReq.FromId, friendReq.TargetId, &exist)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("好友请求已经存在,请等待对方回应")
	}

	//判断是否已经是好友
	err = s.friendMapper.FriendsExist(ctx, friendReq.FromId, friendReq.TargetId, &exist)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("你们已经是好友")
	}
	err = s.friendMapper.CreateFriendReq(ctx, friendReq)
	if err != nil {
		return err
	}
	//推送添加好友申请
	message := models.Message{
		FromId:   friendReq.FromId,
		TargetId: friendReq.TargetId,
		Type:     "friendRequest",
	}

	err = redis.IncrFriendReqUnread(ctx, friendReq.TargetId)
	if err != nil {
		return err
	}
	return s.messageSender.SendWs(&message)

}

func (s *FriendService) RequestList(ctx context.Context, targetId uint) ([]models.FriendApplyResp, error) {
	var list []models.FriendApplyResp
	//查全部，
	err := s.friendMapper.SelectFriendReqListAndInfo(ctx, targetId, &list)
	if err != nil {
		return list, err
	}
	return list, nil
}

// Accept 供控制层调用
func (s *FriendService) Accept(ctx context.Context, fromId uint, targetId uint) error {
	//更新好友状态
	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := s.friendMapper.UpdateStatusWithTx(ctx, tx, fromId, targetId)
		if err != nil {
			return err
		}
		err = s.friendMapper.CreateFriendsWithTx(ctx, tx, fromId, targetId)
		if err != nil {
			return err
		}
		err = s.friendMapper.CreateFriendsWithTx(ctx, tx, targetId, fromId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	key1 := redis.KeyFriendList + strconv.Itoa(int(fromId))
	key2 := redis.KeyFriendList + strconv.Itoa(int(targetId))
	_, err = redis.Rdb.Del(ctx, key1).Result()
	if err != nil {
		fmt.Println("redis error:", err.Error())
	}
	_, err = redis.Rdb.Del(ctx, key2).Result()
	if err != nil {
		fmt.Println("redis error:", err.Error())
	}

	return nil
}

func (s *FriendService) Reject(ctx context.Context, fromId uint, targetId uint) error {
	return s.friendMapper.UpdateStatus(ctx, fromId, targetId)
}

func (s *FriendService) GetFriendList(ctx context.Context, id uint) ([]models.FriendResp, error) {
	var list []models.FriendResp
	//先查redis
	list, err := redis.GetFriendList(ctx, id)
	if err == nil && len(list) > 0 {
		//更新状态，此字段不存redis，单独维护
		for i, friend := range list {
			status, _ := redis.GetUserLine(ctx, friend.Id)
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
		err = redis.SaveFriendList(ctx, id, list)
		if err != nil {
			fmt.Println("更新好友列表缓存失败", err.Error())
		}
	}()

	//更新状态，此字段不存redis，单独维护
	for i, friend := range list {
		status, _ := redis.GetUserLine(ctx, friend.Id)
		list[i].IsOnline = status
	}
	return list, nil
}

func (s *FriendService) UnReadCount(ctx context.Context, userid uint) (int64, error) {
	var count int64
	//先查redis
	unread, err := redis.GetFriendReqUnread(ctx, userid)
	if err == nil {
		return unread, nil
	}
	err = s.friendMapper.FriendReqUnreadCount(userid, &count)
	//写回redis
	go func() {
		if e := recover(); e != nil {
			fmt.Println("好友请求未读数量接口pinic", e)
		}
		err = redis.SetFriendReqUnread(ctx, userid, count)
		if err != nil {
			fmt.Println("更新好友请求未读数量失败", err.Error())
		}
	}()
	if err != nil {
		return count, err
	}
	return count, nil

}

func (s *FriendService) HasRead(ctx context.Context, userId uint) error {
	err := s.friendMapper.FriendReqHasRead(userId)
	if err != nil {
		return err
	}
	_ = redis.SetFriendReqUnread(ctx, userId, 0)
	return nil
}
