package service

import (
	"GinChat/models"
	"GinChat/utils"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

func CreateGroup(userId uint, groupReq *models.CreateGroupReq) error {
	group := models.GroupModel{
		GroupName:  groupReq.GroupName,
		OwnerID:    userId,
		TotalCount: 1,
	}
	err := utils.DB.Transaction(func(tx *gorm.DB) error {
		//创建群聊
		err := tx.Create(&group).Error
		if err != nil {
			return err
		}
		//创建群成员
		err = tx.Create(&models.GroupMember{
			UserID:  userId,
			GroupID: group.ID,
			IsMute:  0,
			Role:    2,
		}).Error
		if err != nil {
			return err
		}
		//创建群主->群的会话
		err = tx.Create(&models.Conversation{
			UserID:      userId,
			PeerID:      group.ID,
			UnreadCount: 0,
			Type:        1,
		}).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil

}

func InviteGroup(inviteReq *models.InviteReq) error {
	err := utils.DB.Transaction(func(tx *gorm.DB) error {
		members := []models.GroupMember{}
		for _, id := range inviteReq.InvitedId {
			members = append(members, models.GroupMember{
				GroupID: inviteReq.GroupId,
				UserID:  id,
			})
		}
		//批量插入
		err := tx.Create(&members).Error
		if err != nil {
			return err
		}
		conversations := []models.Conversation{}
		for _, id := range inviteReq.InvitedId {
			conversations = append(conversations, models.Conversation{
				UserID:      id,
				PeerID:      inviteReq.GroupId,
				UnreadCount: 0,
				Type:        1,
			})
		}

		err = tx.Create(&conversations).Error
		if err != nil {
			return err
		}
		err = tx.Model(&models.GroupModel{}).Where("id=?", inviteReq.GroupId).UpdateColumn("total_count", gorm.Expr("total_count+?", len(inviteReq.InvitedId))).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	//邀请加群成功后，缓存删除
	key := utils.KeyGroupMemberId + strconv.Itoa(int(inviteReq.GroupId))
	go func() {
		if e := recover(); e != nil {
			fmt.Println("邀请入群redisPanic", e)
		}
		_, err2 := utils.Rdb.Del(utils.Ctx, key).Result()
		if err2 != nil {
			fmt.Println("删除群id缓存失败", err2.Error())
		}
	}()
	return nil
}

func GroupDetail(groupId uint64) (models.GroupDetailVO, error) {
	var detail models.GroupDetailVO
	group := models.GroupModel{}
	err := utils.DB.Model(&models.GroupModel{}).Where("id=?", groupId).Take(&group).Error
	if err != nil {
		return detail, err
	}
	members := []models.GroupMemberVO{}
	//只查8个人
	err = utils.DB.Table("group_members gm").Where("gm.group_id=?", groupId).
		Joins("join user_basic u on gm.user_id=u.id").
		Select("u.avatar,u.name,gm.role,gm.user_id as userId").Limit(8).
		Order("gm.role desc").
		Find(&members).Error
	if err != nil {
		return detail, err
	}
	detail = models.GroupDetailVO{
		Avatar:     group.Avatar,
		GroupID:    group.ID,
		GroupName:  group.GroupName,
		TotalCount: group.TotalCount,
		Members:    members,
		Notice:     group.Notice,
	}
	return detail, nil
}

func GroupMembers(groupId uint64, groupMemberReq *models.GroupMemberReq) ([]models.GroupMemberVO, error) {
	var members []models.GroupMemberVO
	err := utils.DB.Debug().Table("group_members gm").Where("gm.group_id=?", groupId).
		Joins("join user_basic u on gm.user_id=u.id").
		Select("u.avatar,u.name,gm.role,gm.user_id as userId").
		Limit(groupMemberReq.PageSize).
		Offset((groupMemberReq.Page - 1) * groupMemberReq.PageSize).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}
