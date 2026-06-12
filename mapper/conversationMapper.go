package mapper

import (
	"GinChat/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ConversationMapper struct {
	db *gorm.DB
}

func NewConversationMapper(db *gorm.DB) *ConversationMapper {
	return &ConversationMapper{
		db: db,
	}
}

func (m *ConversationMapper) ConversationList(userId uint, list *[]models.ConversationInfo) error {
	return m.db.Table("conversations cv").
		Where("cv.user_id = ?", userId).
		Select(
			"cv.id, cv.user_id, cv.peer_id, cv.last_time, cv.last_msg, cv.unread_count, cv.type",
			"CASE WHEN cv.type = 0 THEN u.name WHEN cv.type = 1 THEN g.group_name ELSE '' END AS name",
			"CASE WHEN cv.type = 0 THEN u.avatar WHEN cv.type = 1 THEN g.avatar ELSE '' END AS avatar",
		).
		Joins("LEFT JOIN user_basic u ON cv.type = 0 AND cv.peer_id = u.id").
		Joins("LEFT JOIN group_models g ON cv.type = 1 AND cv.peer_id = g.id").
		Order("cv.last_time DESC").
		Find(list).Error
}

func (m *ConversationMapper) ClearUnreadCount(userId uint, peerId uint64) error {
	return m.db.Model(&models.Conversation{}).Where("user_id = ? and peer_id = ?", userId, peerId).
		UpdateColumn("unread_count", 0).Error
}

func (m *ConversationMapper) UpdateWithTxChat(tx *gorm.DB, message *models.Message, fromId uint, targetId uint) *gorm.DB {
	return tx.Model(&models.Conversation{}).Where("user_id = ? and peer_id = ?", fromId, targetId).
		Updates(map[string]interface{}{
			"last_msg":  message.Content,
			"last_time": message.CreateAt,
		})
}
func (m *ConversationMapper) UpdateWithTxGroup(tx *gorm.DB, message *models.Message) *gorm.DB {
	return tx.Model(&models.Conversation{}).Where("peer_id = ?", message.TargetId).
		UpdateColumns(map[string]interface{}{
			"last_msg":  message.Content,
			"last_time": message.CreateAt,
		})
}
func (m *ConversationMapper) CreateConversationWithTx(tx *gorm.DB, message *models.Message, fromId uint, targetId uint, unread uint) error {
	return tx.Create(&models.Conversation{
		UserID:      fromId,
		PeerID:      targetId,
		LastMsg:     message.Content,
		LastTime:    message.CreateAt,
		UnreadCount: unread,
	}).Error
}
func (m *ConversationMapper) CreateConversationGroupWithTx(tx *gorm.DB, userId uint, groupId uint) error {
	return tx.Create(&models.Conversation{
		UserID:      userId,
		PeerID:      groupId,
		UnreadCount: 0,
		Type:        1,
	}).Error
}
func (m *ConversationMapper) CreateConversationsGroupWithTx(tx *gorm.DB, conversations *[]models.Conversation) error {
	return tx.Create(conversations).Error
}
func (m *ConversationMapper) UpdateUnreadWithTx(tx *gorm.DB, message *models.Message, columns ...clause.Expr) error {
	if len(columns) != 0 {
		return tx.Model(&models.Conversation{}).Where("user_id = ? and peer_id = ?", message.TargetId, message.FromId).
			UpdateColumn("unread_count", columns[0]).Error
	}
	return tx.Model(&models.Conversation{}).Where("user_id = ? and peer_id = ?", message.FromId, message.TargetId).
		UpdateColumn("unread_count", 0).Error
}
