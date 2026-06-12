package mapper

import (
	"GinChat/models"

	"gorm.io/gorm"
)

type GroupMapper struct {
	db *gorm.DB
}

func NewGroupMapper(db *gorm.DB) *GroupMapper {
	return &GroupMapper{
		db: db,
	}
}

func (m *GroupMapper) CreateGroupWithTx(tx *gorm.DB, group *models.GroupModel) error {
	return tx.Create(group).Error
}

func (m *GroupMapper) CreateMemberWithTx(tx *gorm.DB, userId uint, groupId uint) error {
	return tx.Create(&models.GroupMember{
		UserID:  userId,
		GroupID: groupId,
		IsMute:  0,
		Role:    2,
	}).Error
}
func (m *GroupMapper) InviteMemberWithTx(tx *gorm.DB, members *[]models.GroupMember) error {
	return tx.Create(members).Error
}
func (m *GroupMapper) UpdateMemberCountWithTx(tx *gorm.DB, inviteReq *models.InviteReq) error {
	return tx.Model(&models.GroupModel{}).Where("id=?", inviteReq.GroupId).UpdateColumn("total_count", gorm.Expr("total_count+?", len(inviteReq.InvitedId))).Error
}
func (m *GroupMapper) GetGroupInfo(groupId uint64, group *models.GroupModel) error {
	return m.db.Model(&models.GroupModel{}).Where("id=?", groupId).Take(group).Error
}
func (m *GroupMapper) GetMember8Info(groupId uint64, members *[]models.GroupMemberVO) error {
	return m.db.Table("group_members gm").Where("gm.group_id=?", groupId).
		Joins("join user_basic u on gm.user_id=u.id").
		Select("u.avatar,u.name,gm.role,gm.user_id as userId").Limit(8).
		Order("gm.role desc").
		Find(members).Error
}
func (m *GroupMapper) GetMemberPageInfo(groupId uint64, groupMemberReq *models.GroupMemberReq, members *[]models.GroupMemberVO) error {
	return m.db.Table("group_members gm").Where("gm.group_id=?", groupId).
		Joins("join user_basic u on gm.user_id=u.id").
		Select("u.avatar,u.name,gm.role,gm.user_id as userId").
		Limit(groupMemberReq.PageSize).
		Offset((groupMemberReq.Page - 1) * groupMemberReq.PageSize).Find(members).Error
}
func (m *GroupMapper) GetAllMemberId(targetId uint, memberIds *[]uint) error {
	return m.db.Model(&models.GroupMember{}).Select("user_id").Where("group_id=?", targetId).
		Find(memberIds).Error
}
