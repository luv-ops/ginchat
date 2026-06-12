package service

import (
	"GinChat/mapper"
	"GinChat/models"
)

type ConversationService struct {
	conversationMapper *mapper.ConversationMapper
}

func NewConversationService(cM *mapper.ConversationMapper) *ConversationService {
	return &ConversationService{
		conversationMapper: cM,
	}
}
func (s *ConversationService) ConversationList(userId uint) ([]models.ConversationInfo, error) {
	list := []models.ConversationInfo{}
	err := s.conversationMapper.ConversationList(userId, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *ConversationService) ClearUnreadCount(userId uint, peerId uint64) error {
	return s.conversationMapper.ClearUnreadCount(userId, peerId)
}
