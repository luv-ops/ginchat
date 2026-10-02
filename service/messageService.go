package service

import (
	"GinChat/mapper"
	"GinChat/models"
	"context"

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
func (s *MessageService) GetMessage(ctx context.Context, userId uint, messageReq *models.MessageReq) ([]models.MessageVO, error) {
	var list []models.MessageVO
	var db *gorm.DB
	// 单聊历史记录
	if messageReq.Type == "chat" {
		db = s.messageMapper.ChatMessage(ctx, userId, messageReq)
	} else {
		// 群聊历史记录
		db = s.messageMapper.ChatGroup(ctx, messageReq)
	}
	err := s.messageMapper.MessagePage(ctx, db, messageReq, &list)

	if err != nil {
		return nil, err
	}
	return list, nil

}

//Where("from_id = ? and target_id = ? or from_id = ? and target_id = ?", userId, messageReq.PeerId, messageReq.PeerId, userId)
