package service

import (
	"GinChat/mapper"
	"GinChat/models"

	"gorm.io/gorm"
)

type MessageService struct {
	messageMapper *mapper.MessageMapper
}

func NewMessageService(mM *mapper.MessageMapper) *MessageService {
	return &MessageService{
		messageMapper: mM,
	}
}
func (s *MessageService) GetMessage(userId uint, messageReq *models.MessageReq) ([]models.MessageVO, error) {
	var list []models.MessageVO
	var db *gorm.DB
	// 单聊历史记录
	if messageReq.Type == "chat" {
		db = s.messageMapper.ChatMessage(userId, messageReq)
	} else {
		// 群聊历史记录
		db = s.messageMapper.ChatGroup(messageReq)
	}
	err := s.messageMapper.MessagePage(db, messageReq, &list)

	if err != nil {
		return nil, err
	}
	return list, nil

}

//Where("from_id = ? and target_id = ? or from_id = ? and target_id = ?", userId, messageReq.PeerId, messageReq.PeerId, userId)
