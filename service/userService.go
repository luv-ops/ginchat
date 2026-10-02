package service

import (
	"GinChat/mapper"
	"GinChat/models"
	"GinChat/utils"
	"context"
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
func (s *UserService) GetUserList(ctx context.Context) (*[]models.UserBasic, error) {
	var data []models.UserBasic
	err := s.userMapper.GetUserList(ctx, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *UserService) Register(ctx context.Context, body *models.RegisterReq) error {

	var exist bool
	err := s.userMapper.UserExistByName(ctx, body.Name, &exist)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("用户已经存在")
	}

	encode, err := utils.Encode(body.Password)
	if err != nil {
		return err
	}
	user := models.UserBasic{Name: body.Name, Password: encode}
	err = s.userMapper.CreateUser(ctx, &user)
	//应对高并发
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errors.New("用户已经存在")
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Login(ctx context.Context, body *models.LoginReq) (*models.UserBasic, error) {

	user := models.UserBasic{}
	//根据名字查用户
	err := s.userMapper.SelectByName(ctx, body.Name, &user)
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
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.userMapper.DeleteById(ctx, id)

}

func (s *UserService) UpdateUser(ctx context.Context, body *models.UpdateReq, id uint) error {
	var exist bool
	err := s.userMapper.UserExistById(ctx, id, &exist)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("用户不存在")
	}
	return s.userMapper.UpdateById(ctx, id, body)
}

func (s *UserService) UserInfo(ctx context.Context, userId uint) (models.UserBasic, error) {
	user := models.UserBasic{}
	err := s.userMapper.GetUserInfoById(userId, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}
