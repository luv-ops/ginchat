package service

import (
	"GinChat/Mysql"
	"GinChat/models"
	"GinChat/redis"
	"fmt"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
)

type wsMessage struct {
	Mt      int
	Message models.Message
}

var (
	wsConnMap = make(map[uint]*websocket.Conn) //后期可以用redis 存储
	wsLock    sync.RWMutex                     //读写锁,读可以共享，写独占
)

// WsConnectionHandler
// 管理用户map连接
func WsConnectionHandler(connect *websocket.Conn, userId uint) {

	//加入map在线表
	wsLock.Lock()
	wsConnMap[userId] = connect
	redis.SetUserOnline(userId, 1)
	wsLock.Unlock()
	fmt.Println("用户上线：", userId)
	//清理资源,断开连接
	defer func() {
		//删除用户连接
		wsLock.Lock()
		delete(wsConnMap, userId)
		redis.SetUserOnline(userId, 0)
		wsLock.Unlock()
		connect.Close()
	}()
	//定义管道，for循环开始解耦，分为读监听和写监听
	//为什么需要在for循环加go？因为for循环，读消息是阻塞任务，不能被阻塞,所以引入go进行异步并发
	messageChan := make(chan wsMessage, 100)
	// 退出信号：任一协程断开，发送信号
	exitChan := make(chan struct{})
	go readLoop(connect, messageChan, exitChan)
	go writeLoop(messageChan, userId, exitChan)
	//select阻塞，不让函数结束，因为函数退出，连接会关闭，不能为空select，会永久阻塞，连接不会断开
	select {
	case <-exitChan:
		//收到退出信号，关闭连接，执行defer
		return
	}
}

// 读监听
func readLoop(connect *websocket.Conn, messageChan chan wsMessage, exitChan chan struct{}) {
	defer func() {
		close(exitChan)
	}()
	for {
		// 读取前端消息
		mt, message, err := connect.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println("接收到消息：", string(message))
		msg := models.Message{}
		json.Unmarshal(message, &msg)
		switch msg.Type {
		case "chat":
			Mysql.DB.Create(&msg)
		}
		tempMessage := wsMessage{
			Mt:      mt,
			Message: msg,
		}
		messageChan <- tempMessage
	}
}

// 写监听
func writeLoop(messageChan chan wsMessage, userId uint, exitChan chan struct{}) {

	for {
		select {
		case <-exitChan:
			return
		case obj := <-messageChan:
			obj = <-messageChan
			msg := obj.Message
			msg.FromId = userId
			//将消息序列化为字节数组
			data, _ := json.Marshal(msg)
			wsLock.RLock()
			conn := wsConnMap[msg.TargetId]
			wsLock.RUnlock()
			//目标不在线
			if conn == nil {
				continue
			}
			err := conn.WriteMessage(obj.Mt, data)
			if err != nil {
				fmt.Println("发送消息失败")
			}
		}

	}
}

// SendWs  单向
// 供服务器其他业务层发送ws
func SendWs(msg *models.Message) {
	// 读锁，支持并发，性能更高,RLock读锁专用，性能更高
	wsLock.RLock()
	ws := wsConnMap[msg.TargetId]
	wsLock.RUnlock()
	if ws == nil {
		return
	}
	var data []byte
	switch msg.Type {
	case "friendRequest":
		user := models.UserBasic{}
		Mysql.DB.Take(&user, msg.FromId)

		reqMsg := models.FriendApplyResp{
			FromId:   msg.FromId,
			Avatar:   user.Avatar,
			Type:     msg.Type,
			CreateAt: time.Now().Format(time.DateTime),
			Msg:      "",
			Name:     user.Name,
			Status:   0,
		}
		data, _ = json.Marshal(reqMsg)
	case "chat":
		data, _ = json.Marshal(msg)
	}

	go func() {
		ws.WriteMessage(1, data)
	}()

	//这里为什么不用select阻塞函数,因为函数退出，连接不会关闭
}

// SendWsGroup  群聊
func SendWsGroup(msg *models.MessageVO) {
	//先获取群成员所有id
	var memberIds []uint
	ids, err := redis.GetGroupMemberIds(msg.TargetId)

	if err == nil && len(ids) > 0 {
		memberIds = ids
	} else {
		err = Mysql.DB.Model(&models.GroupMember{}).Select("user_id").Where("group_id=?", msg.TargetId).
			Find(&memberIds).Error
		if err != nil {
			fmt.Println("获取群成员id失败", err.Error())
			return
		}
		//缓存构建
		redis.SetGroupMemberIds(msg.TargetId, memberIds)
	}
	//遍历所有成员
	for _, id := range memberIds {
		//跳过发送者
		if id == msg.FromId {
			continue
		}
		wsLock.RLock()
		ws := wsConnMap[id]
		wsLock.RUnlock()
		if ws == nil {
			continue
		}
		go func() {
			data, _ := json.Marshal(msg)
			ws.WriteMessage(1, data)
		}()

	}
}
