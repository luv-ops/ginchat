package service

import (
	"GinChat/mapper"
	"GinChat/models"
	"GinChat/utils"
	"errors"

	"gorm.io/gorm"
)

type UserService struct {
	userMapper *mapper.UserMapper
}

func NewUserService(uM *mapper.UserMapper) *UserService {
	return &UserService{
		userMapper: uM,
	}
}
func (s *UserService) GetUserList() (*[]models.UserBasic, error) {
	var data []models.UserBasic
	err := s.userMapper.GetUserList(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *UserService) Register(body *models.RegisterReq) error {

	var exist bool
	err := s.userMapper.UserExistByName(body.Name, &exist)
	//为了用户体验
	if exist {
		return errors.New("用户已经存在")
	}
	//密码加密
	encode, err := utils.Encode(body.Password)
	if err != nil {
		return err
	}
	user := models.UserBasic{Name: body.Name, Password: encode}
	err = s.userMapper.CreateUser(&user)
	//应对高并发
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errors.New("用户已经存在")
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Login(body *models.LoginReq) (*models.UserBasic, error) {

	user := models.UserBasic{}
	//根据名字查用户
	err := s.userMapper.SelectByName(body.Name, &user)
	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("用户不存在")
	}
	if !utils.Verify(user.Password, body.Password) {
		return nil, errors.New("密码错误")
	}
	return &user, nil
}
func (s *UserService) DeleteUser(id uint) error {
	return s.userMapper.DeleteById(id)

}

func (s *UserService) UpdateUser(body *models.UpdateReq, id uint) error {
	var exist bool
	err := s.userMapper.UserExistById(id, &exist)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("用户不存在")
	}
	return s.userMapper.UpdateById(id, body)
}

func (s *UserService) UserInfo(userId uint) (models.UserBasic, error) {
	user := models.UserBasic{}
	err := s.userMapper.GetUserInfoById(userId, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}
