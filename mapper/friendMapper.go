package mapper

import (
	"GinChat/models"

	"gorm.io/gorm"
)

type FriendMapper struct {
	db *gorm.DB
}

func NewFriendMapper(db *gorm.DB) *FriendMapper {
	return &FriendMapper{
		db: db,
	}
}

func (m *FriendMapper) FriendReqExist(fromId uint, targetId uint, count *int64) error {
	return m.db.Model(&models.FriendReq{}).Where("from_id = ? and target_id = ?", fromId, targetId).
		Where("status = ?", 0).Count(count).Error
}
func (m *FriendMapper) FriendsExist(fromId uint, targetId uint, friendCount *int64) error {
	return m.db.Model(&models.Friends{}).Where("user_id = ? and friend_id = ?", fromId, targetId).
		Where("status = ?", 1).Count(friendCount).Error
}
func (m *FriendMapper) CreateFriendReq(friendReq *models.FriendReq) error {
	return m.db.Create(friendReq).Error
}
func (m *FriendMapper) SelectFriendReqListAndInfo(targetId uint, list *[]models.FriendApplyResp) error {
	return m.db.Table("friend_reqs fq").Where("from_id = ? or target_id = ?", targetId, targetId).
		Where("status in ?", []string{"0", "3"}).
		//如果发起者是我，就拿对方 ID 去关联名字
		//如果发起者是对方，就拿对方 ID 去关联名字
		Joins("LEFT JOIN user_basic ON "+
			"CASE WHEN fq.from_id = ? THEN fq.target_id ELSE fq.from_id END = user_basic.id", targetId).
		Select("fq.from_id,fq.target_id,fq.status,user_basic.name,user_basic.avatar,fq.create_at").
		Order("fq.create_at desc").
		Find(list).Error

}
func (m *FriendMapper) SelectFriendListAndInfo(id uint, list *[]models.FriendResp) error {
	return m.db.Model(&models.Friends{}).Joins("left join user_basic on user_basic.id = friends.friend_id").
		Select("user_basic.id,user_basic.name,user_basic.avatar").
		Find(list, "friends.user_id = ?", id).Error

}
func (m *FriendMapper) UpdateStatusWithTx(tx *gorm.DB, fromId uint, targetId uint) error {
	return tx.Model(&models.FriendReq{}).
		Where("from_id = ? and target_id = ? and status in ?", fromId, targetId, []string{"0", "3"}). //只处理状态为0的请求
		UpdateColumn("status", 1).Error
}
func (m *FriendMapper) CreateFriendsWithTx(tx *gorm.DB, fromId uint, targetId uint) error {
	return tx.Create(&models.Friends{
		UserId:   fromId,
		FriendId: targetId,
	}).Error
}
func (m *FriendMapper) UpdateStatus(fromId uint, targetId uint) error {
	return m.db.Model(&models.FriendReq{}).
		Where("from_id = ? and target_id = ? and status in ?", fromId, targetId, []string{"0", "3"}). //只处理状态为0的请求
		UpdateColumn("status", 1).Error
}
func (m *FriendMapper) FriendReqUnreadCount(userid uint, count *int64) error {
	return m.db.Model(&models.FriendReq{}).Where("target_id = ? and status = ?", userid, 0).
		Count(count).Error
}
func (m *FriendMapper) FriendReqHasRead(userId uint) error {
	return m.db.Model(&models.FriendReq{}).Where("target_id = ? and status = ?", userId, 0).
		Update("status", 3).Error
}
