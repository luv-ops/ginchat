package service

import (
	"GinChat/Mysql"
	"GinChat/models"
	"GinChat/redis"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

func AddFriend(friendReq *models.FriendReq) error {
	var count int64
	//查询用户是否存在
	err := Mysql.DB.Take(&models.UserBasic{}, friendReq.TargetId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}
	Mysql.DB.Model(&models.FriendReq{}).Where("from_id = ? and target_id = ?", friendReq.FromId, friendReq.TargetId).
		Where("status = ?", 0).Count(&count)
	if count > 0 {
		return errors.New("好友请求已经存在,请等待对方回应")
	}
	//判断是否已经是好友
	var friendCount int64
	Mysql.DB.Model(&models.Friends{}).Where("user_id = ? and friend_id = ?", friendReq.FromId, friendReq.TargetId).
		Where("status = ?", 1).Count(&friendCount)
	if friendCount > 0 {
		return errors.New("你们已经是好友")
	}
	err = Mysql.DB.Create(friendReq).Error
	if err != nil {
		return err
	}
	user := &models.UserBasic{}
	Mysql.DB.Take(&user, friendReq.FromId)
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

	SendWs(&message)
	return nil
}

func RequestList(targetId uint) ([]models.FriendApplyResp, error) {
	list := []models.FriendApplyResp{}
	//查全部，
	err := Mysql.DB.Table("friend_reqs fq").Where("from_id = ? or target_id = ?", targetId, targetId).
		Where("status in ?", []string{"0", "3"}).
		//如果发起者是我，就拿对方 ID 去关联名字
		//如果发起者是对方，就拿对方 ID 去关联名字
		Joins("LEFT JOIN user_basic ON "+
			"CASE WHEN fq.from_id = ? THEN fq.target_id ELSE fq.from_id END = user_basic.id", targetId).
		Select("fq.from_id,fq.target_id,fq.status,user_basic.name,user_basic.avatar,fq.create_at").
		Order("fq.create_at desc").
		Find(&list).Error

	if err != nil {
		return list, err
	}
	for e := range list {
		fmt.Println(list[e])
	}
	return list, nil
}

func Accept(fromId uint, targetId uint) error {
	err := Mysql.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.FriendReq{}).
			Where("from_id = ? and target_id = ? and status in ?", fromId, targetId, []string{"0", "3"}). //只处理状态为0的请求
			UpdateColumn("status", 1).Error
		if err != nil {
			return err
		}
		err = tx.Create(&models.Friends{
			UserId:   fromId,
			FriendId: targetId,
		}).Error
		if err != nil {
			return err
		}
		err = tx.Create(&models.Friends{
			UserId:   targetId,
			FriendId: fromId,
		}).Error
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
func Reject(fromId uint, targetId uint) error {
	err := Mysql.DB.Model(&models.FriendReq{}).
		Where("from_id = ? and target_id = ? and status in ?", fromId, targetId, []string{"0", "3"}).
		UpdateColumn("status", 2).Error
	if err != nil {
		return err
	}
	return nil
}

func GetFriendList(id uint) ([]models.FriendResp, error) {
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
	err = Mysql.DB.Model(&models.Friends{}).Joins("left join user_basic on user_basic.id = friends.friend_id").
		Select("user_basic.id,user_basic.name,user_basic.avatar").
		Find(&list, "friends.user_id = ?", id).Error

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

func UnReadCount(userid uint) (int64, error) {
	var count int64
	//先查redis
	unread, err := redis.GetFriendReqUnread(userid)
	if err == nil {
		return unread, nil
	}
	err = Mysql.DB.Model(&models.FriendReq{}).Where("target_id = ? and status = ?", userid, 0).
		Count(&count).Error
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

func HasRead(userId uint) error {
	err := Mysql.DB.Model(&models.FriendReq{}).Where("target_id = ? and status = ?", userId, 0).
		Update("status", 3).Error
	if err != nil {
		return err
	}
	redis.SetFriendReqUnread(userId, 0)
	return nil

}
