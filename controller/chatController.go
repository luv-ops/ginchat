package controller

import (
	"GinChat/models"
	"GinChat/service"
	"GinChat/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	chatService *service.ChatService
}
type IChatController interface {
	Send(c *gin.Context)
}

func NewChatController(cs *service.ChatService) *ChatController {
	return &ChatController{
		chatService: cs,
	}
}

// Send
// @Tags 聊天模块
// @Summary 发送消息
// @Param data body models.Message true "聊天参数"
// @Success 200 {object} utils.Response{data=EmptyData}
// @Router /upload/send [post]
func (con *ChatController) Send(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "请先登录")
		return
	}
	message := models.Message{}
	err := c.ShouldBindJSON(&message)
	if userId == message.TargetId {
		utils.Fail(c, http.StatusBadRequest, "不能给自己发消息")
		return
	}
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if message.Content == "" {
		utils.Fail(c, http.StatusBadRequest, "消息不能为空")
		return
	}
	message.FromId = userId.(uint)
	err = con.chatService.Send(c.Request.Context(), &message)

	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Ok2(c, "发送成功")
}
