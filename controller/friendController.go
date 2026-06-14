package controller

import (
	"GinChat/models"
	"GinChat/service"
	"GinChat/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FriendController struct {
	friendService *service.FriendService
}
type IFriendController interface {
	AddFriend(c *gin.Context)
	RequestList(c *gin.Context)
	UnReadCount(c *gin.Context)
	Accept(c *gin.Context)
	Reject(c *gin.Context)
	GetFriendList(c *gin.Context)
	HasRead(c *gin.Context)
}

func NewFriendController(fs *service.FriendService) *FriendController {
	return &FriendController{
		friendService: fs,
	}
}

// AddFriend
// @Tags 好友模块
// @Summary 添加好友
// @Param data body models.FriendReq true "登录参数"
// @Success 200 {object} utils.Response{}
// @Router /friend/add [post]
func (con *FriendController) AddFriend(c *gin.Context) {

	fromId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	friendReq := models.FriendReq{}
	if fromId == friendReq.FromId {
		utils.Fail(c, 400, "不能添加自己为好友")
		return
	}
	err := c.ShouldBindJSON(&friendReq)
	friendReq.FromId = fromId.(uint)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	err = con.friendService.AddFriend(c.Request.Context(), &friendReq)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Ok2(c, "添加好友成功")
}

// RequestList
// @Tags 好友模块
// @Summary 获取好友请求列表
// @Success 200 {object} utils.Response{data=[]models.FriendApplyResp}
// @Router /friend/requests [get]
func (con *FriendController) RequestList(c *gin.Context) {
	targetId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	list, err := con.friendService.RequestList(targetId.(uint))
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.Ok(c, list)
}

// UnReadCount
// @Tags 好友模块
// @Summary 获取好友请求未读计数
// @Success 200 {object} utils.Response{data=map[string]int64}
// @Router /friend/requests/unread [get]
func (con *FriendController) UnReadCount(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	count, err := con.friendService.UnReadCount(userId.(uint))
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.Ok(c, map[string]int64{"unread": count})
}

// Accept
// @Tags 好友模块
// @Summary 同意好友申请
// @Param friendId path string true "好友ID"
// @Success 200 {object} utils.Response{}
// @Router /friend/accept/{friendId} [post]
func (con *FriendController) Accept(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	friendId := c.Param("friendId")
	fmt.Println(friendId)
	fid, err := strconv.ParseUint(friendId, 10, 64)
	if friendId == "" || err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}
	err = con.friendService.Accept(uint(fid), userId.(uint))
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.Ok2(c, "添加好友成功")

}

// Reject
// @Tags 好友模块
// @Summary 拒绝好友申请
// @Param friendId path string true "好友ID"
// @Success 200 {object} utils.Response{}
// @Router /friend/reject/{friendId} [post]
func (con *FriendController) Reject(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	friendId := c.Param("friendId")
	fmt.Println(friendId)
	fid, err := strconv.ParseUint(friendId, 10, 64)
	if friendId == "" || err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}
	err = con.friendService.Reject(uint(fid), userId.(uint))
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.Ok2(c, "拒绝好友申请成功")
}

// GetFriendList
// @Tags 好友模块
// @Summary 获取好友列表
// @Success 200 {object} utils.Response{data=models.FriendResp}
// @Router /friend/list [get]
func (con *FriendController) GetFriendList(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	list, err := con.friendService.GetFriendList(userId.(uint))
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.Ok(c, list)
}

func (con *FriendController) HasRead(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	err := con.friendService.HasRead(userId.(uint))
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.Ok2(c, "已读成功")
}
