package controller

import (
	"GinChat/models"
	"GinChat/service"
	"GinChat/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMessage
// @Tags 消息模块
// @Summary 获取消息历史记录
// @Param peerId query uint true "对方用户ID"
// @Param page query int false "页码，默认1"
// @Param size query int false "每页条数，默认20"
// @Success 200 {object} utils.Response{data=[]models.Message}
// @Router /message/list [get]
func GetMessage(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	messageReq := models.MessageReq{}
	err := c.ShouldBindQuery(&messageReq)

	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if messageReq.Size == 0 {
		messageReq.Size = 20
	}
	message, err := service.GetMessage(userId.(uint), &messageReq)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, "获取消息失败")
		return
	}
	utils.Ok(c, message)
}
