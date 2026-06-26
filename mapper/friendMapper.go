package mapper

import (
	"GinChat/models"

	"gorm.io/gorm"
)

type FriendMapper struct {
	db *gorm.DB
}

func NewFriendMapper(db *gorm.DB) *FriendMapper {
	return &FriendMapper{
		db: db,
	}
}

// FriendReqExist 检查两个用户之间是否存在好友请求

func (m *FriendMapper) FriendReqExist(fromId uint, targetId uint, exist *bool) error {

	return m.db.Raw("select exists(select 1 from friend_reqs where from_id = ? and target_id = ? and status = ?)", fromId, targetId, 0).Scan(exist).Error
}

// FriendsExist 检查两个用户之间是否存在好友关系

func (m *FriendMapper) FriendsExist(fromId uint, targetId uint, friendExist *bool) error {

	return m.db.Raw("select exists(select 1 from friends where user_id = ? and friend_id = ?)", fromId, targetId).Scan(friendExist).Error
}

// CreateFriendReq 方法用于创建好友请求记录

func (m *FriendMapper) CreateFriendReq(friendReq *models.FriendReq) error {
	return m.db.Create(friendReq).Error
}

// SelectFriendReqListAndInfo 查询好友请求列表及相关信息

func (m *FriendMapper) SelectFriendReqListAndInfo(targetId uint, list *[]models.FriendApplyResp) error {
	// 执行数据库查询操作
	return m.db.Table("friend_reqs fq").Where("from_id = ? or target_id = ?", targetId, targetId).
		Where("status in ?", []string{"0", "3"}).
		//如果发起者是我，就拿对方 ID 去关联名字
		//如果发起者是对方，就拿对方 ID 去关联名字
		Joins("LEFT JOIN user_basic ON "+
			"CASE WHEN fq.from_id = ? THEN fq.target_id ELSE fq.from_id END = user_basic.id", targetId).
		Select("fq.from_id,fq.target_id,fq.status,user_basic.name,user_basic.avatar,fq.create_at").
		Order("fq.create_at desc").
		Find(list).Error

}

// SelectFriendListAndInfo 根据用户ID查询其好友列表及基本信息

func (m *FriendMapper) SelectFriendListAndInfo(id uint, list *[]models.FriendResp) error {
	// 使用数据库模型进行查询
	// 1. 从Friends表开始
	// 2. 左连接UserBasic表，关联条件为friends.friend_id = user_basic.id
	// 3. 选择UserBasic表中的id、name和avatar字段
	// 4. 查询条件为friends.user_id等于传入的id
	// 5. 将查询结果存入list指针
	return m.db.Model(&models.Friends{}).Joins("left join user_basic on user_basic.id = friends.friend_id").
		Select("user_basic.id,user_basic.name,user_basic.avatar").
		Find(list, "friends.user_id = ?", id).Error

}

// UpdateStatusWithTx 使用事务更新好友请求状态的方法

func (m *FriendMapper) UpdateStatusWithTx(tx *gorm.DB, fromId uint, targetId uint) error {
	return tx.Model(&models.FriendReq{}). // 指定操作的数据模型为FriendReq
						Where("from_id = ? and target_id = ? and status in ?", fromId, targetId, []string{"0", "3"}). //只处理状态为0的请求
						UpdateColumn("status", 1).Error
}

// CreateFriendsWithTx 在事务中创建好友关系记录

func (m *FriendMapper) CreateFriendsWithTx(tx *gorm.DB, fromId uint, targetId uint) error {
	// 使用事务创建好友关系记录
	// 创建一个Friends结构体实例并保存到数据库
	return tx.Create(&models.Friends{
		UserId:   fromId,   // 用户ID
		FriendId: targetId, // 好友ID
	}).Error
}

// UpdateStatus 更新好友请求的状态

func (m *FriendMapper) UpdateStatus(fromId uint, targetId uint) error {
	return m.db.Model(&models.FriendReq{}). // 使用数据库模型FriendReq进行操作
						Where("from_id = ? and target_id = ? and status in ?", fromId, targetId, []string{"0", "3"}). //只处理状态为0的请求
						UpdateColumn("status", 1).Error
}

// FriendReqUnreadCount 是一个方法，用于获取用户未读的好友请求数量

func (m *FriendMapper) FriendReqUnreadCount(userid uint, count *int64) error {
	// 使用数据库模型查询目标用户ID且状态为0(未读)的好友请求数量
	// 将查询结果存储到count指针指向的变量中
	// 如果操作过程中出现错误，则返回错误信息
	return m.db.Model(&models.FriendReq{}).Where("target_id = ? and status = ?", userid, 0).
		Count(count).Error
}

// FriendReqHasRead 是一个方法，用于将指定用户的所有未读好友请求标记为已读

func (m *FriendMapper) FriendReqHasRead(userId uint) error {
	// 使用数据库模型进行操作
	// 1. 查找目标ID为userId且状态为0(未读)的所有好友请求
	// 2. 将这些记录的状态更新为3(已读)
	// 3. 如果操作过程中出现错误，返回错误信息
	return m.db.Model(&models.FriendReq{}).Where("target_id = ? and status = ?", userId, 0).
		Update("status", 3).Error
}
