package service

import (
	"GinChat/mapper"
	"GinChat/models"
	"GinChat/redis"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type GroupService struct {
	groupMapper        *mapper.GroupMapper
	conversationMapper *mapper.ConversationMapper
	db                 *gorm.DB
}

func NewGroupService(gM *mapper.GroupMapper, cM *mapper.ConversationMapper, db *gorm.DB) *GroupService {
	return &GroupService{
		groupMapper:        gM,
		conversationMapper: cM,
		db:                 db,
	}
}
func (s *GroupService) CreateGroup(userId uint, groupReq *models.CreateGroupReq) error {
	group := models.GroupModel{
		GroupName:  groupReq.GroupName,
		OwnerID:    userId,
		TotalCount: 1,
	}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		//创建群聊
		err := s.groupMapper.CreateGroupWithTx(tx, &group)
		if err != nil {
			return err
		}
		//创建群成员
		err = s.groupMapper.CreateMemberWithTx(tx, userId, group.ID)
		if err != nil {
			return err
		}
		//创建群主->群的会话
		err = s.conversationMapper.CreateConversationGroupWithTx(tx, userId, group.ID)
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

func (s *GroupService) InviteGroup(inviteReq *models.InviteReq) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		members := []models.GroupMember{}
		for _, id := range inviteReq.InvitedId {
			members = append(members, models.GroupMember{
				GroupID: inviteReq.GroupId,
				UserID:  id,
			})
		}
		//批量插入
		err := s.groupMapper.InviteMemberWithTx(tx, &members)
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

		err = s.conversationMapper.CreateConversationsGroupWithTx(tx, &conversations)
		if err != nil {
			return err
		}
		err = s.groupMapper.UpdateMemberCountWithTx(tx, inviteReq)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	//邀请加群成功后，缓存删除
	key := redis.KeyGroupMemberId + strconv.Itoa(int(inviteReq.GroupId))
	go func() {
		if e := recover(); e != nil {
			fmt.Println("邀请入群redisPanic", e)
		}
		_, err2 := redis.Rdb.Del(redis.Ctx, key).Result()
		if err2 != nil {
			fmt.Println("删除群id缓存失败", err2.Error())
		}
	}()
	return nil
}

func (s *GroupService) GroupDetail(groupId uint64) (models.GroupDetailVO, error) {
	var detail models.GroupDetailVO
	group := models.GroupModel{}
	err := s.groupMapper.GetGroupInfo(groupId, &group)
	if err != nil {
		return detail, err
	}
	members := []models.GroupMemberVO{}
	//只查8个人
	err = s.groupMapper.GetMember8Info(groupId, &members)
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

func (s *GroupService) GroupMembers(groupId uint64, groupMemberReq *models.GroupMemberReq) ([]models.GroupMemberVO, error) {
	var members []models.GroupMemberVO
	err := s.groupMapper.GetMemberPageInfo(groupId, groupMemberReq, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}
