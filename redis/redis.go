package redis

import (
	"GinChat/models"
	"errors"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
)

const (
	KeyUserOnline      = "user:online:"
	KeyFriendReqUnread = "friend:request:unread:"
	UnreadExpire       = 3 * 24 * time.Hour
	KeyGroupMemberId   = "group:member:"
	KeyFriendList      = "friend:list:"
	UserOffline        = 0
	LockSession        = "lock:session:"
)

func SetUserOnline(userId uint, status int) error {
	uid := strconv.Itoa(int(userId))
	key := KeyUserOnline + uid
	// 先清空旧数据（避免重复）
	Rdb.Del(Ctx, key)
	return Rdb.Set(Ctx, key, status, 7*24*time.Hour).Err()
}

func GetUserLine(userId uint) (int, error) {
	id := strconv.Itoa(int(userId))
	key := KeyUserOnline + id
	val, err := Rdb.Get(Ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		//key不存在默认离线
		return UserOffline, nil
	}
	return val, nil
}
func SetFriendReqUnread(userId uint, count int64) error {
	id := strconv.Itoa(int(userId))
	key := KeyFriendReqUnread + id
	// 先清空旧数据（避免重复）
	Rdb.Del(Ctx, key)
	Rdb.Expire(Ctx, key, UnreadExpire)
	return Rdb.Set(Ctx, key, count, UnreadExpire).Err()
}
func GetFriendReqUnread(userId uint) (int64, error) {
	id := strconv.Itoa(int(userId))
	val, err := Rdb.Get(Ctx, KeyFriendReqUnread+id).Int64()
	return val, err
}
func IncrFriendReqUnread(userId uint) error {
	id := strconv.Itoa(int(userId))
	_, err := Rdb.Incr(Ctx, KeyFriendReqUnread+id).Result()
	if err != nil {
		return err
	}
	//续期
	Rdb.Expire(Ctx, KeyFriendReqUnread+id, UnreadExpire)
	return nil
}

func SetGroupMemberIds(groupId uint, memberIds []uint) error {
	key := KeyGroupMemberId + strconv.Itoa(int(groupId))
	if len(memberIds) == 0 {
		return nil
	}
	// 先清空旧数据（避免重复）
	Rdb.Del(Ctx, key)
	//将memberids全部id转string
	var members []interface{}
	for _, id := range memberIds {
		members = append(members, strconv.Itoa(int(id)))
	}
	Rdb.Expire(Ctx, key, 5*24*time.Hour)
	return Rdb.SAdd(Ctx, key, members...).Err()
}
func GetGroupMemberIds(groupId uint) ([]uint, error) {
	key := KeyGroupMemberId + strconv.Itoa(int(groupId))
	members, err := Rdb.SMembers(Ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var ids []uint
	for _, memberIds := range members {
		num, _ := strconv.ParseUint(memberIds, 10, 64)
		ids = append(ids, uint(num))
	}
	//Redis 的 SMembers 行为：
	//key 不存在 → 返回空切片 [] + nil error
	//key 存在但为空 → 返回空切片 [] + nil error
	return ids, nil
}

func SaveFriendList(userId uint, friends []models.FriendResp) error {
	id := strconv.Itoa(int(userId))
	if len(friends) == 0 {
		return nil
	}
	key := KeyFriendList + id
	// 先清空旧数据（避免重复）
	Rdb.Del(Ctx, key)
	pipe := Rdb.Pipeline()
	for _, friend := range friends {
		data, err := json.Marshal(friend)
		if err != nil {
			return err
		}
		pipe.LPush(Ctx, key, data)
	}
	pipe.Expire(Ctx, key, 3*24*time.Hour)
	_, err := pipe.Exec(Ctx)
	return err
}
func GetFriendList(userId uint) ([]models.FriendResp, error) {
	id := strconv.Itoa(int(userId))
	key := KeyFriendList + id
	var friends []models.FriendResp
	result, err := Rdb.LRange(Ctx, key, 0, -1).Result()
	if err != nil {
		return friends, err
	}
	for _, friend := range result {
		var friendResp models.FriendResp
		err = json.Unmarshal([]byte(friend), &friendResp)
		if err != nil {
			return friends, err
		}
		friends = append(friends, friendResp)
	}
	return friends, nil
}
