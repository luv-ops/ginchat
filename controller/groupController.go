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

type GroupController struct {
	groupService *service.GroupService
}
type IGroupController interface {
	CreateGroup(c *gin.Context)
	InviteGroup(c *gin.Context)
	GroupDetail(c *gin.Context)
	GroupMembers(c *gin.Context)
	TestJoin(c *gin.Context)
}

func NewGroupController(gS *service.GroupService) *GroupController {
	return &GroupController{
		groupService: gS,
	}
}

// CreateGroup
// @Tags 群组模块
// @Summary 创建群聊
// @Param data body models.CreateGroupReq true "创建群参数"
// @Success 200 {object} utils.Response{data=EmptyData}
// @Router /group/create [post]
func (con *GroupController) CreateGroup(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	var groupReq models.CreateGroupReq
	err := c.ShouldBindJSON(&groupReq)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if groupReq.GroupName == "" {
		utils.Fail(c, http.StatusBadRequest, "群名称不能为空")
		return
	}
	err = con.groupService.CreateGroup(c.Request.Context(), userId.(uint), &groupReq)
	if err != nil {
		utils.Fail(c, 500, "创建群聊失败")
		return
	}
	utils.Ok2(c, "创建群聊成功")
}

// InviteGroup
// @Tags 群组模块
// @Summary 邀请加入群
// @Param data body models.InviteReq true "邀请入群参数"
// @Success 200 {object} utils.Response{data=EmptyData}
// @Router /group/invite [post]
func (con *GroupController) InviteGroup(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	inviteReq := models.InviteReq{}
	err := c.ShouldBindJSON(&inviteReq)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	err = con.groupService.InviteGroup(c.Request.Context(), &inviteReq)
	if err != nil {
		utils.Fail(c, 500, "邀请失败")
		return
	}
	utils.Ok2(c, "邀请成功")
}

// GroupDetail
// @Tags 群组模块
// @Summary 获取群详情
// @Param groupId path uint true "群组ID"
// @Success 200 {object} utils.Response{data=models.GroupDetailVO}
// @Router /group/detail/{groupId} [get]
func (con *GroupController) GroupDetail(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
	}
	var groupId uint64
	tempId := c.Param("groupId")
	groupId, err := strconv.ParseUint(tempId, 10, 64)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	detail, err := con.groupService.GroupDetail(groupId)
	if err != nil {
		utils.Fail(c, 500, "获取群详情失败")
		return
	}
	utils.Ok(c, detail)
}

// GroupMembers
// @Tags 群组模块
// @Summary 获取群详情
// @Param groupId path uint true "群组ID"
// @Param data query models.GroupMemberReq true "查询参数"
// @Success 200 {object} utils.Response{data=models.GroupMemberVO}
// @Router /group/members/{groupId} [get]
func (con *GroupController) GroupMembers(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "用户未登录")
	}
	var groupId uint64
	var groupMemberReq models.GroupMemberReq
	tempId := c.Param("groupId")
	groupId, err := strconv.ParseUint(tempId, 10, 64)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	err = c.ShouldBindQuery(&groupMemberReq)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	fmt.Println(groupMemberReq)
	members, err := con.groupService.GroupMembers(groupId, &groupMemberReq)
	if err != nil {
		utils.Fail(c, 500, "获取群成员失败")
		return
	}
	utils.Ok(c, members)
}

func (con *GroupController) TestJoin(c *gin.Context) {
	var groupId uint = 5
	con.groupService.JoinGroup(groupId)

}
