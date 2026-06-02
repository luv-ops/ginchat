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

// AddFriend
// @Tags 好友模块
// @Summary 添加好友
// @Param data body models.FriendReq true "登录参数"
// @Success 200 {object} utils.Response{}
// @Router /friend/add [post]
func AddFriend(c *gin.Context) {

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
	err = service.AddFriend(&friendReq)
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
func RequestList(c *gin.Context) {
	targetId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	list, err := service.RequestList(targetId.(uint))
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
func UnReadCount(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	count, err := service.UnReadCount(userId.(uint))
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
func Accept(c *gin.Context) {
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
	err = service.Accept(uint(fid), userId.(uint))
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
func Reject(c *gin.Context) {
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
	err = service.Reject(uint(fid), userId.(uint))
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
func GetFriendList(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	list, err := service.GetFriendList(userId.(uint))
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.Ok(c, list)
}

func HasRead(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	err := service.HasRead(userId.(uint))
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.Ok2(c, "已读成功")
}
