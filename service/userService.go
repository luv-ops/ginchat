package service

import (
	"GinChat/MQ"
	"GinChat/mapper"
	"GinChat/models"
	"GinChat/utils"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserService struct {
	userMapper *mapper.UserMapper
	kafkaCli   *MQ.KafkaClient
}

func NewUserService(uM *mapper.UserMapper, kc *MQ.KafkaClient) *UserService {
	return &UserService{
		userMapper: uM,
		kafkaCli:   kc,
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

func (s *UserService) Register(ctx context.Context, body *models.RegisterReq) error {

	var exist bool
	err := s.userMapper.UserExistByName(body.Name, &exist)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("用户已经存在")
	}
	dto := MQ.UserDTO{
		Name:     body.Name,
		Password: body.Password,
	}
	return s.kafkaCli.ProduceUser(ctx, &dto, MQ.TopicUserCreate)
}
func (s *UserService) HandleUserCreate(dto *MQ.UserDTO) error {
	password := dto.Password
	name := dto.Name
	//密码加密
	encode, err := utils.Encode(password)
	if err != nil {
		return err
	}
	user := models.UserBasic{Name: name, Password: encode}
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
