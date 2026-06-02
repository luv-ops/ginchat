package controller

import (
	"GinChat/service"
	"GinChat/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 定义 Upgrader（只在这里定义）
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
}

// ConnectWs
// @Summary WebSocket 连接websocket接口
// @Description 建立 WebSocket 长连接，用于实时聊天
// @Tags 聊天模块
// @Accept json
// @Produce json
// @Param token path string true "用户认证token"
// @Success 101 {string} string "WebSocket 连接成功"
// @Failure 400 {object} string "连接失败"
// @Router /ws/{token} [get]
func ConnectWs(c *gin.Context) {
	//鉴权
	userId, ok := c.Get("userId")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "请先登录")
		return
	}
	//建立ws连接,协议升级
	connect, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("ws连接失败", err)
		return
	}

	//处理ws连接,加入/删除用户连接map
	service.WsConnectionHandler(connect, userId.(uint))

}
