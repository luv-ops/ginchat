package service

import (
	"GinChat/models"
	"GinChat/utils"
)

func ConversationList(userId uint) ([]models.ConversationInfo, error) {
	list := []models.ConversationInfo{}
	err := utils.DB.Table("conversations cv").
		Where("cv.user_id = ?", userId).
		Select(
			"cv.id, cv.user_id, cv.peer_id, cv.last_time, cv.last_msg, cv.unread_count, cv.type",
			"CASE WHEN cv.type = 0 THEN u.name WHEN cv.type = 1 THEN g.group_name ELSE '' END AS name",
			"CASE WHEN cv.type = 0 THEN u.avatar WHEN cv.type = 1 THEN g.avatar ELSE '' END AS avatar",
		).
		Joins("LEFT JOIN user_basic u ON cv.type = 0 AND cv.peer_id = u.id").
		Joins("LEFT JOIN group_models g ON cv.type = 1 AND cv.peer_id = g.id").
		Order("cv.last_time DESC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func ClearUnreadCount(userId uint, peerId uint64) error {
	err := utils.DB.Model(&models.Conversation{}).Where("user_id = ? and peer_id = ?", userId, peerId).
		UpdateColumn("unread_count", 0).Error
	if err != nil {
		return err
	}
	return nil
}
