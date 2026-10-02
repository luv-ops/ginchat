package service

import (
	"GinChat/mapper"
	"GinChat/models"
	"GinChat/redis"
	"context"
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
func (s *GroupService) CreateGroup(ctx context.Context, userId uint, groupReq *models.CreateGroupReq) error {

	group := models.GroupModel{
		GroupName:  groupReq.GroupName,
		OwnerID:    userId,
		TotalCount: 1,
	}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		//创建群聊
		err := s.groupMapper.CreateGroupWithTx(ctx, tx, &group)
		if err != nil {
			return err
		}
		//创建群成员
		err = s.groupMapper.CreateMemberWithTx(ctx, tx, userId, group.ID)
		if err != nil {
			return err
		}
		//创建群主->群的会话
		err = s.conversationMapper.CreateConversationGroupWithTx(ctx, tx, userId, group.ID)
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

func (s *GroupService) InviteGroup(ctx context.Context, inviteReq *models.InviteReq) error {
	//获取哪些id已经在群里面了

	var existIds []uint
	err := s.groupMapper.ExistsMemberIds(ctx, &inviteReq.InvitedId, inviteReq.GroupId, &existIds)
	if err != nil {
		return err
	}
	existMap := make(map[uint]bool)
	for _, id := range existIds {
		existMap[id] = true
	}
	var needAddIds []uint
	for _, id := range inviteReq.InvitedId {
		if !existMap[id] {
			needAddIds = append(needAddIds, id)
		}
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		var members []models.GroupMember
		for _, id := range inviteReq.InvitedId {
			members = append(members, models.GroupMember{
				GroupID: inviteReq.GroupId,
				UserID:  id,
			})
		}
		//批量插入
		err1 := s.groupMapper.InviteMemberWithTx(ctx, tx, &members)
		if err1 != nil {
			return err1
		}
		var conversations []models.Conversation
		for _, id := range inviteReq.InvitedId {
			conversations = append(conversations, models.Conversation{
				UserID:      id,
				PeerID:      inviteReq.GroupId,
				UnreadCount: 0,
				Type:        1,
			})
		}

		err = s.conversationMapper.CreateConversationsGroupWithTx(ctx, tx, &conversations)
		if err != nil {
			return err
		}
		err = s.groupMapper.UpdateMemberCountWithTx(ctx, tx, inviteReq)
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

	_, err2 := redis.Rdb.Del(ctx, key).Result()
	if err2 != nil {
		fmt.Println("删除群id缓存失败", err2.Error())
	}

	return nil

}
func (s *GroupService) GroupDetail(ctx context.Context, groupId uint64) (models.GroupDetailVO, error) {
	var detail models.GroupDetailVO
	group := models.GroupModel{}
	err := s.groupMapper.GetGroupInfo(ctx, groupId, &group)
	if err != nil {
		return detail, err
	}
	var members []models.GroupMemberVO
	//只查8个人
	err = s.groupMapper.GetMember8Info(ctx, groupId, &members)
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

func (s *GroupService) GroupMembers(ctx context.Context, groupId uint64, groupMemberReq *models.GroupMemberReq) ([]models.GroupMemberVO, error) {
	var members []models.GroupMemberVO
	err := s.groupMapper.GetMemberPageInfo(ctx, groupId, groupMemberReq, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (s *GroupService) JoinGroup(ctx context.Context, groupId uint) error {
	var groupMembers []models.GroupMember
	for i := 7805; i <= 9804; i++ {
		groupMembers = append(groupMembers, models.GroupMember{
			GroupID: groupId,
			UserID:  uint(i),
		})
	}
	return s.groupMapper.JoinGroup(ctx, &groupMembers)
}
