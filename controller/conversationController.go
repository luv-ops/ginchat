package controller

import (
	"GinChat/service"
	"GinChat/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConversationController struct {
	conversationService *service.ConversationService
}
type IConversationController interface {
	ConversationList(c *gin.Context)
	ClearUnreadCount(c *gin.Context)
}

func NewConversationController(cs *service.ConversationService) *ConversationController {
	return &ConversationController{
		conversationService: cs,
	}
}

// ConversationList
// @Tags 会话模块
// @Summary 获取会话列表
// @Success 200 {object} utils.Response{data=[]models.ConversationInfo}
// @Router /conversation/list [get]
func (con *ConversationController) ConversationList(c *gin.Context) {
	userId, ok := c.Get("userId")

	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	list, err := con.conversationService.ConversationList(userId.(uint))
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, "获取会话列表失败")
		return
	}
	utils.Ok(c, list)
}

// ClearUnreadCount
// @Tags 会话模块
// @Summary 清楚当前会话的未读计数
// @Param peerId path string true "对方用户ID"
// @Success 200 {object} utils.Response{}
// @Router /conversation/unreadClear/{peerId} [post]
func (con *ConversationController) ClearUnreadCount(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
	}
	pId := c.Param("peerId")
	var peerId uint64
	peerId, err := strconv.ParseUint(pId, 10, 64)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	err = con.conversationService.ClearUnreadCount(userId.(uint), peerId)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, "清除未读计数失败")
		return
	}
	utils.Ok2(c, "清除未读计数成功")
}
