package service

import (
	"GinChat/Mysql"
	"GinChat/models"
	"GinChat/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func GetUserList() (*[]models.UserBasic, error) {
	var data []models.UserBasic
	err := Mysql.DB.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func Register(body *models.RegisterReq) error {

	var count int64
	err := Mysql.DB.Model(&models.UserBasic{}).Where("name=?", body.Name).Count(&count).Error
	fmt.Println("count", count)
	//为了用户体验
	if count > 0 {
		return errors.New("用户已经存在")
	}
	//密码加密
	encode, slat := utils.Md5Encode(body.Password)
	user := models.UserBasic{Name: body.Name, Password: encode, Salt: slat}
	err = Mysql.DB.Create(&user).Error
	//应对高并发
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errors.New("用户已经存在")
	}
	if err != nil {
		return err
	}
	return nil
}

func Login(body *models.LoginReq) (*models.UserBasic, error) {

	user := models.UserBasic{}
	//根据名字查用户
	err := Mysql.DB.Take(&user, "name=?", body.Name).Error
	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("用户不存在")
	}
	if !utils.Equal(body.Password, user.Salt, user.Password) {
		return nil, errors.New("密码错误")
	}
	return &user, nil
}
func DeleteUser(id uint) error {
	err := Mysql.DB.Delete(&models.UserBasic{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(body *models.UpdateReq, id uint) error {
	err := Mysql.DB.Model(&models.UserBasic{ID: id}).UpdateColumns(body).Error
	if err != nil {
		return err
	}
	return nil
}

func UserInfo(userId uint) (models.UserBasic, error) {
	user := models.UserBasic{}
	err := Mysql.DB.Take(&user, userId).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
