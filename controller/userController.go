package controller

import (
	"GinChat/models"
	"GinChat/service"
	"GinChat/utils"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator/v12"
	"github.com/gin-gonic/gin"
)

type EmptyData struct{}

// GetUserList
// @Tags 用户模块
// @Summary 获取用户列表
// @Success 200 {object} utils.Response{data=[]models.UserRes}
// @Router /user/list [get]
func GetUserList(c *gin.Context) {

	data, err := service.GetUserList()
	if err != nil {
		utils.Fail(c, 500, "删除失败")
		return
	}
	utils.Ok(c, data)
}

// Login
// @Tags 用户模块
// @Summary 登录
// @Param data body models.LoginReq true "登录参数"
// @Success 200 {object} utils.Response{data=models.UserRes}
// @Router /user/login [post]
func Login(c *gin.Context) {
	body := models.LoginReq{}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	user, err := service.Login(&body)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		utils.Fail(c, 500, "token生成失败")
		return
	}
	res := models.UserRes{
		Id:     user.ID,
		Avatar: user.Avatar,
		Name:   user.Name,
		Token:  "Bearer " + token,
	}
	utils.Ok(c, res)
}

// Register
// @Tags 用户模块
// @Summary 注册
// @Param data body models.RegisterReq true "登录参数"
// @Success 200 {object} utils.Response{data=EmptyData}
// @Router /user/register [post]
func Register(c *gin.Context) {
	var body models.RegisterReq
	err := c.ShouldBindJSON(&body)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	if body.Password != body.ConfirmPassword {
		utils.Fail(c, 400, "前后密码不一致")
		return
	}
	err = service.Register(&body)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}
	utils.Ok2(c, "创建用户成功")
}

// DeleteUser
// @Tags 用户模块
// @Summary 删除用户
// @Success 200 {object} utils.Response{data=EmptyData}
// @Router /user/delete [post]
func DeleteUser(c *gin.Context) {
	id, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	//类型断言
	err := service.DeleteUser(id.(uint))
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}
	utils.Ok2(c, "删除用户成功")
}

// UpdateUser
// @Tags 用户模块
// @Summary 修改用户消息
// @Param data body models.UpdateReq true "修改参数"
// @Success 200 {object} utils.Response{data=EmptyData}
// @Router /user/update [post]
func UpdateUser(c *gin.Context) {
	body := models.UpdateReq{}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	if body.Email != "" && !govalidator.IsEmail(body.Email) {
		utils.Fail(c, 400, "邮箱格式错误")
		return
	}
	if body.Phone != "" && !utils.IsPhone(body.Phone) {
		utils.Fail(c, 400, "手机格式错误")
		return
	}
	id, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	err = service.UpdateUser(&body, id.(uint))
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}
	utils.Ok2(c, "更新成功")
}

// UserInfo
// @Tags 用户模块
// @Summary 获取用户信息
// @Param userId query int true "用户ID"
// @Success 200 {object} utils.Response{data=[]models.UserRes}
// @Router /user/info [get]
func UserInfo(c *gin.Context) {
	id := c.Query("userId")
	if id == "" {
		utils.Fail(c, 400, "参数错误")
		return
	}
	userId, _ := strconv.ParseUint(id, 10, 64)
	userInfo, err := service.UserInfo(uint(userId))
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}
	utils.Ok(c, userInfo)
}
