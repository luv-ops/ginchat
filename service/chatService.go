package service

import (
	"GinChat/MQ"
	"GinChat/mapper"
	"GinChat/models"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatService struct {
	userMapper         *mapper.UserMapper
	conversationMapper *mapper.ConversationMapper
	messageMapper      *mapper.MessageMapper
	messageSender      IMessageSender
	db                 *gorm.DB
	kafkaCli           *MQ.KafkaClient
}

func NewChatService(uM *mapper.UserMapper, cM *mapper.ConversationMapper, mM *mapper.MessageMapper,
	mS IMessageSender, db *gorm.DB, kC *MQ.KafkaClient) *ChatService {
	return &ChatService{
		userMapper:         uM,
		conversationMapper: cM,
		messageMapper:      mM,
		messageSender:      mS,
		db:                 db,
		kafkaCli:           kC,
	}
}

// 供controller层使用
func (s *ChatService) Send(ctx context.Context, message *models.Message) error {
	if message.Content == "" {
		return errors.New("消息内容不能为空")
	}
	var mType int
	switch message.Type {
	case "chat":
		mType = MQ.ChatTypePrivate
	case "groupMessage":
		mType = MQ.ChatTypeGroup
	default:
		return errors.New("未知的消息类型")
	}
	//组装kafka消息
	dto := MQ.MsgDTO{
		MsgID:    uuid.NewString(),
		FromID:   message.FromId,
		TargetID: message.TargetId,
		ChatType: mType,
		MsgType:  message.MsgType,
		Content:  message.Content,
	}
	var err error
	start := time.Now()
	if mType == MQ.ChatTypePrivate {
		err = s.kafkaCli.SendPrivateMsg(ctx, &dto)
	} else {
		err = s.kafkaCli.SendGroupMsg(ctx, &dto)
	}
	cost := time.Since(start)
	log.Printf("【Kafka发送耗时】%s", cost)

	return err
}

// 供kafka异步调用，并通过实现防止循环引用
func (s *ChatService) HandleMsg(dto *MQ.MsgDTO) error {
	var mType string
	switch dto.ChatType {
	case MQ.ChatTypePrivate:
		mType = "chat"
	case MQ.ChatTypeGroup:
		mType = "groupMessage"
	}
	message := &models.Message{
		FromId:   dto.FromID,
		TargetId: dto.TargetID,
		MsgType:  dto.MsgType,
		Type:     mType,
		Content:  dto.Content,
	}
	err := s.messageMapper.CreateMessage(message)
	if err != nil {
		return err
	}
	//更新会话没有会话则创建
	err = s.updateConversation(message)
	if err != nil {
		return err
	}
	//查询fromId的用户名和头像
	user := models.UserBasic{}
	err = s.userMapper.GetUserInfoById(message.FromId, &user, "name,avatar")
	if err != nil {
		return err
	}
	messageVO := models.MessageVO{
		FromName:   user.Name,
		FromAvatar: user.Avatar,
		Message:    *message,
	}
	switch message.Type {
	case "chat":
		fmt.Println("单聊消息", message)
		return s.messageSender.SendWs(message)
	case "groupMessage":
		fmt.Println("群聊消息", message)
		return s.messageSender.SendWsGroup(&messageVO)
	}
	return nil
}

func (s *ChatService) updateConversation(message *models.Message) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		//尝试更新，如果RowsAffected == 0则会话不存在则创建
		//自己发的消息，UpdateColumn不走钩子，性能比update高  同UpdateColumns比updates
		//更新我们的
		err1 := s.sender(tx, message)
		var err2 error
		//只有单聊才需要双向
		if message.Type == "chat" {
			//更新对方的
			err2 = s.receiver(tx, message)
		}
		if err1 != nil || err2 != nil {
			if err1 != nil {
				return err1
			}
			return err2
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *ChatService) sender(tx *gorm.DB, message *models.Message) error {
	fmt.Println("发送方", message)
	var res1 *gorm.DB
	switch message.Type {
	case "chat":

		res1 = s.conversationMapper.UpdateWithTxChat(tx, message, message.FromId, message.TargetId)
	case "groupMessage":
		res1 = s.conversationMapper.UpdateWithTxGroup(tx, message)
	default:
		return errors.New("未知的消息类型")
	}
	if res1.Error != nil {
		return res1.Error
	}
	//说明会话记录不存在
	if res1.RowsAffected == 0 {
		err := s.conversationMapper.CreateConversationWithTx(tx, message, message.FromId, message.TargetId, 0)
		if err != nil {
			return err
		}
	} else {
		//会话存在更新未读计数
		err := s.conversationMapper.UpdateUnreadWithTx(tx, message)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *ChatService) receiver(tx *gorm.DB, message *models.Message) error {
	res := s.conversationMapper.UpdateWithTxChat(tx, message, message.TargetId, message.FromId)
	if res.Error != nil {
		return res.Error
	}
	//说明会话记录不存在
	if res.RowsAffected == 0 {
		err := s.conversationMapper.CreateConversationWithTx(tx, message, message.TargetId, message.FromId, 1)
		if err != nil {
			return err
		}
	} else {
		//会话存在更新未读计数
		err := s.conversationMapper.UpdateUnreadWithTx(tx, message, gorm.Expr("unread_count + ?", 1))
		if err != nil {
			return err
		}
	}
	return nil
}
