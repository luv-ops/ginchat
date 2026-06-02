package service

import (
	"GinChat/models"
	"GinChat/utils"
)

func GetMessage(userId uint, messageReq *models.MessageReq) ([]models.MessageVO, error) {
	list := []models.MessageVO{}
	db := utils.DB.Model(&models.Message{})
	// 单聊历史记录
	if messageReq.Type == "chat" {
		db = db.Where("(from_id = ? AND target_id = ? AND type = ?) OR (from_id = ? AND target_id = ? AND type = ?)",
			userId, messageReq.PeerId, messageReq.Type,
			messageReq.PeerId, userId, messageReq.Type,
		)
	} else {
		// 群聊历史记录
		db = db.Table("messages ms").
			Select("ms.id, ms.from_id, ms.target_id, ms.type, ms.content, ms.media, ms.picture, ms.url, ms.create_at, "+
				"ub.name as from_name, ub.avatar as from_avatar").
			Joins("join user_basic ub on ms.from_id = ub.id").
			Where("target_id = ? AND type = ?", messageReq.PeerId, messageReq.Type)
	}
	err := db.Order("create_at DESC").
		Limit(messageReq.Size).Offset(messageReq.Size * (messageReq.Page - 1)).
		Scan(&list).Error

	if err != nil {
		return nil, err
	}
	return list, nil

}

//Where("from_id = ? and target_id = ? or from_id = ? and target_id = ?", userId, messageReq.PeerId, messageReq.PeerId, userId)
