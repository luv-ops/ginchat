package mapper

import (
	"GinChat/models"

	"gorm.io/gorm"
)

type MessageMapper struct {
	db *gorm.DB
}

func NewMessageMapper(db *gorm.DB) *MessageMapper {
	return &MessageMapper{
		db,
	}
}

func (m *MessageMapper) ChatMessage(userId uint, messageReq *models.MessageReq) *gorm.DB {
	return m.db.Model(&models.Message{}).Where("(from_id = ? AND target_id = ? AND type = ?) OR (from_id = ? AND target_id = ? AND type = ?)",
		userId, messageReq.PeerId, messageReq.Type,
		messageReq.PeerId, userId, messageReq.Type,
	)
}
func (m *MessageMapper) ChatGroup(messageReq *models.MessageReq) *gorm.DB {
	return m.db.Model(&models.Message{}).Table("messages ms").
		Select("ms.id, ms.from_id, ms.target_id, ms.type, ms.msg_type,ms.content, ms.media, ms.picture, ms.url, ms.create_at, "+
			"ub.name as from_name, ub.avatar as from_avatar").
		Joins("join user_basic ub on ms.from_id = ub.id").
		Where("target_id = ? AND type = ?", messageReq.PeerId, messageReq.Type)
}

func (m *MessageMapper) MessagePage(newDb *gorm.DB, messageReq *models.MessageReq, list *[]models.MessageVO) error {
	return newDb.Order("create_at DESC").
		Limit(messageReq.Size).Offset(messageReq.Size * (messageReq.Page - 1)).
		Scan(list).Error
}
func (m *MessageMapper) CreateMessage(message *models.Message) error {
	return m.db.Create(message).Error
}
