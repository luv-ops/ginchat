package mapper

import (
	"GinChat/models"

	"gorm.io/gorm"
)

type UserMapper struct {
	db *gorm.DB
}

func NewUserMapper(db *gorm.DB) *UserMapper {
	return &UserMapper{
		db: db,
	}
}

func (m *UserMapper) GetUserList(data *[]models.UserBasic) error {
	return m.db.Find(data).Error

}
func (m *UserMapper) CreateUser(user *models.UserBasic) error {
	return m.db.Create(user).Error
}
func (m *UserMapper) UserExistByName(name string, exist *bool) error {
	return m.db.Raw("select exists(select 1 from user_basic where name = ?)", name).Scan(exist).Error
}
func (m *UserMapper) UserExistById(id uint, exist *bool) error {
	return m.db.Raw("select exists(select 1 from user_basic where id = ?)", id).Scan(exist).Error
}
func (m *UserMapper) SelectByName(name string, user *models.UserBasic) error {
	return m.db.Take(user, "name=?", name).Error
}

func (m *UserMapper) DeleteById(id uint) error {
	return m.db.Delete(&models.UserBasic{}, id).Error
}

func (m *UserMapper) UpdateById(id uint, body *models.UpdateReq) error {
	return m.db.Model(&models.UserBasic{ID: id}).UpdateColumns(body).Error
}
func (m *UserMapper) GetUserInfoById(userId uint, user *models.UserBasic, columns ...string) error {
	if len(columns) != 0 {
		return m.db.Select(columns).Take(user, userId).Error
	}
	return m.db.Take(user, userId).Error
}
