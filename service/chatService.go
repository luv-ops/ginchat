package service

import (
	"GinChat/models"
	"GinChat/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func Send(message *models.Message) error {
	err := utils.DB.Create(message).Error
	if err != nil {
		return err
	}
	//更新会话没有会话则创建
	err = updateConversation(message)
	if err != nil {
		return err
	}
	//查询fromId的用户名和头像
	user := models.UserBasic{}
	err = utils.DB.Model(&user).Select("name,avatar").Where("id = ?", message.FromId).Take(&user).Error
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
		SendWs(message)
	case "groupMessage":
		fmt.Println("群聊消息", message)
		SendWsGroup(&messageVO)
	}
	return nil
}

func updateConversation(message *models.Message) error {
	fmt.Printf("%+v", message)

	err := utils.DB.Transaction(func(tx *gorm.DB) error {
		//尝试更新，如果RowsAffected == 0则会话不存在则创建
		//自己发的消息，UpdateColumn不走钩子，性能比update高  同UpdateColumns比updates
		//更新我们的
		err1 := sender(tx, message)
		var err2 error
		//只有单聊才需要双向
		if message.Type == "chat" {
			//更新对方的
			err2 = receiver(tx, message)
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

func sender(tx *gorm.DB, message *models.Message) error {
	fmt.Println("发送方", message)
	var res1 *gorm.DB
	switch message.Type {
	case "chat":
		res1 = tx.Model(&models.Conversation{}).Where("user_id = ? and peer_id = ?", message.FromId, message.TargetId).
			Updates(map[string]interface{}{
				"last_msg":  message.Content,
				"last_time": message.CreateAt,
			})
	case "groupMessage":
		res1 = tx.Model(&models.Conversation{}).Where("peer_id = ?", message.TargetId).
			UpdateColumns(map[string]interface{}{
				"last_msg":  message.Content,
				"last_time": message.CreateAt,
			})
	default:
		return errors.New("未知的消息类型")
	}
	if res1.Error != nil {
		return res1.Error
	}
	//说明会话记录不存在
	if res1.RowsAffected == 0 {
		err := tx.Create(&models.Conversation{
			UserID:      message.FromId,
			PeerID:      message.TargetId,
			LastMsg:     message.Content,
			LastTime:    message.CreateAt,
			UnreadCount: 0,
		}).Error
		if err != nil {
			return err
		}
	} else {
		//会话存在更新未读计数
		err := tx.Model(&models.Conversation{}).Where("user_id = ? and peer_id = ?", message.FromId, message.TargetId).
			UpdateColumn("unread_count", 0).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func receiver(tx *gorm.DB, message *models.Message) error {
	res := tx.Model(&models.Conversation{}).Where("user_id = ? and peer_id = ?", message.TargetId, message.FromId).
		Updates(map[string]interface{}{
			"last_msg":  message.Content,
			"last_time": message.CreateAt,
		})
	if res.Error != nil {
		return res.Error
	}
	//说明会话记录不存在
	if res.RowsAffected == 0 {
		err := tx.Create(&models.Conversation{
			UserID:      message.TargetId,
			PeerID:      message.FromId,
			LastMsg:     message.Content,
			LastTime:    message.CreateAt,
			UnreadCount: 1,
		}).Error
		if err != nil {
			return err
		}
	} else {
		//会话存在更新未读计数
		err := tx.Model(&models.Conversation{}).Where("user_id = ? and peer_id = ?", message.TargetId, message.FromId).
			UpdateColumn("unread_count", gorm.Expr("unread_count + ?", 1)).Error
		if err != nil {
			return err
		}
	}
	return nil
}
