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
		db: db, // 将传入的数据库连接对象赋值给GroupMapper的db字段
	}
}

// CreateGroupWithTx 使用事务创建群组的方法

func (m *GroupMapper) CreateGroupWithTx(tx *gorm.DB, group *models.GroupModel) error {
	return tx.Create(group).Error
}

// CreateMemberWithTx 在事务中创建群组成员记录

func (m *GroupMapper) CreateMemberWithTx(tx *gorm.DB, userId uint, groupId uint) error {
	// 使用事务创建群组成员记录
	// 设置初始状态为未禁言(IsMute=0)
	// 设置初始角色为普通成员(Role=2)
	return tx.Create(&models.GroupMember{
		UserID:  userId,  // 用户ID
		GroupID: groupId, // 群组ID
		IsMute:  0,       // 是否禁用，0表示未禁用
		Role:    2,       // 成员角色，2表示普通成员
	}).Error
}

// InviteMemberWithTx 在事务中邀请新成员加入群组

func (m *GroupMapper) InviteMemberWithTx(tx *gorm.DB, members *[]models.GroupMember) error {
	// 使用事务创建新的群组成员记录
	// 如果操作成功，返回nil；如果失败，返回相应的错误信息
	return tx.Create(members).Error
}

// UpdateMemberCountWithTx 使用事务更新群组成员数量

func (m *GroupMapper) UpdateMemberCountWithTx(tx *gorm.DB, inviteReq *models.InviteReq) error {
	// 使用GORM的Model方法指定要更新的表(models.GroupModel{})
	// 使用Where条件筛选出特定的群组(通过inviteReq.GroupId)
	// 使用UpdateColumn更新total_count字段
	// 通过gorm.Expr实现SQL表达式"total_count+?"，其中?被替换为len(inviteReq.InvitedId)
	// 即将群组成员总数增加被邀请的用户数量
	return tx.Model(&models.GroupModel{}).Where("id=?", inviteReq.GroupId).UpdateColumn("total_count", gorm.Expr("total_count+?", len(inviteReq.InvitedId))).Error
}

// GetGroupInfo 根据群组ID获取群组信息

func (m *GroupMapper) GetGroupInfo(groupId uint64, group *models.GroupModel) error {
	// 使用GORM的Model方法指定操作模型为GroupModel
	// Where方法添加查询条件，根据id字段匹配groupId
	// Take方法查询第一条符合条件的结果并填充到group指针中
	// 如果查询失败或未找到记录，返回Error
	return m.db.Model(&models.GroupModel{}).Where("id=?", groupId).Take(group).Error
}

// GetMember8Info 获取指定群组的前8位成员信息

func (m *GroupMapper) GetMember8Info(groupId uint64, members *[]models.GroupMemberVO) error {
	// 执行数据库查询，获取群组成员信息
	return m.db.Table("group_members gm").Where("gm.group_id=?", groupId).
		// 关联查询用户基本信息表
		Joins("join user_basic u on gm.user_id=u.id").
		Select("u.avatar,u.name,gm.role,gm.user_id as userId").Limit(8).
		Order("gm.role desc").
		Find(members).Error
}

// GetMemberPageInfo 获取群组成员分页信息

func (m *GroupMapper) GetMemberPageInfo(groupId uint64, groupMemberReq *models.GroupMemberReq, members *[]models.GroupMemberVO) error {
	// 执行数据库查询，获取群组成员的分页信息
	return m.db.Table("group_members gm").Where("gm.group_id=?", groupId).
		// 关联查询用户基本信息表
		Joins("join user_basic u on gm.user_id=u.id").
		// 选择需要的字段：用户头像、用户名、群成员角色和用户ID
		Select("u.avatar,u.name,gm.role,gm.user_id as userId").
		// 设置查询结果数量限制
		Limit(groupMemberReq.PageSize).
		// 计算偏移量，实现分页
		Offset((groupMemberReq.Page - 1) * groupMemberReq.PageSize).Find(members).Error
}

// GetAllMemberId 获取指定群组中所有成员的用户ID

func (m *GroupMapper) GetAllMemberId(targetId uint, memberIds *[]uint) error {

	return m.db.Model(&models.GroupMember{}).Select("user_id").Where("group_id=?", targetId).
		Find(memberIds).Error
}

// MemberExistsGroup 检查成员是否存在于组中的方法

func (m *GroupMapper) MemberExistsGroup(userId uint, groupId uint, memberExist *bool) error {
	return m.db.Raw("select exists(select 1 from group_members where user_id = ? and group_id = ?)", userId, groupId).Scan(memberExist).Error
}

func (m *GroupMapper) ExistsMemberIds(inviteIds *[]uint, groupId uint, existIds *[]uint) error {
	return m.db.Model(&models.GroupMember{}).Where(" user_id in (?) and group_id = ?", &inviteIds, groupId).Pluck("user_id", existIds).Error
}
