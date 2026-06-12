package service

import (
	"GinChat/mapper"
	"GinChat/models"
	"GinChat/redis"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
)

type WebsocketService struct {
	groupMapper *mapper.GroupMapper
	userMapper  *mapper.UserMapper
}

func NewWebsocketService(gM *mapper.GroupMapper, uM *mapper.UserMapper) *WebsocketService {
	return &WebsocketService{
		groupMapper: gM,
		userMapper:  uM,
	}
}

type IMessageSender interface {
	SendWs(msg *models.Message) error
	SendWsGroup(msg *models.MessageVO) error
}

var (
	wsConnMap = make(map[uint]*websocket.Conn) //后期可以用redis 存储
	wsLock    sync.RWMutex                     //读写锁,读可以共享，写独占
)

// WsConnectionHandler
// 管理用户map连接
func (s *WebsocketService) WsConnectionHandler(connect *websocket.Conn, userId uint) {

	//加入map在线表
	wsLock.Lock()
	wsConnMap[userId] = connect
	_ = redis.SetUserOnline(userId, 1)
	wsLock.Unlock()
	fmt.Println("用户上线：", userId)
	//清理资源,断开连接
	defer func() {
		//删除用户连接
		wsLock.Lock()
		delete(wsConnMap, userId)
		_ = redis.SetUserOnline(userId, 0)
		wsLock.Unlock()
		_ = connect.Close()
	}()
	//为什么需要在for循环加go？因为for循环，读消息是阻塞任务，不能被阻塞,所以引入go进行异步并发
	// 退出信号：任一协程断开，发送信号
	exitChan := make(chan struct{})
	go readLoop(connect, userId, exitChan)
	//select阻塞，不让函数结束，因为函数退出，连接会关闭，不能为空select，会永久阻塞，连接不会断开
	select {
	case <-exitChan:
		//收到退出信号，关闭连接，执行defer
		return
	}
}

// 读监听
func readLoop(connect *websocket.Conn, userId uint, exitChan chan struct{}) {
	defer func() {
		close(exitChan)
	}()
	for {
		_, message, err := connect.ReadMessage()
		if err != nil {
			fmt.Printf("用户 %d 断开连接: %v\n", userId, err)
			break
		}

		// 可以处理心跳消息
		msg := string(message)
		if msg == "ping" || msg == "heartbeat" {
			// 回复心跳
			_ = connect.WriteMessage(websocket.TextMessage, []byte("pong"))
		}

	}
}

// SendWs  单向
// 供服务器其他业务层发送ws
func (s *WebsocketService) SendWs(msg *models.Message) error {
	// 读锁，支持并发，性能更高,RLock读锁专用，性能更高
	wsLock.RLock()
	ws := wsConnMap[msg.TargetId]
	wsLock.RUnlock()
	if ws == nil {
		return errors.New("用户不在线")
	}
	var data []byte
	switch msg.Type {
	case "friendRequest":
		user := models.UserBasic{}
		err := s.userMapper.GetUserInfoById(msg.FromId, &user)
		if err != nil {
			return err
		}
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
		if e := recover(); e != nil {
			fmt.Println("发送ws错误", e)
		}
		err := ws.WriteMessage(1, data)
		if err != nil {
			fmt.Println("发送ws错误", err)
		}
	}()
	return nil
}

// SendWsGroup  群聊
func (s *WebsocketService) SendWsGroup(msg *models.MessageVO) error {
	//先获取群成员所有id
	var memberIds []uint
	ids, err := redis.GetGroupMemberIds(msg.TargetId)
	if err != nil {
		fmt.Println("缓存获取群成员id失败", err.Error())
	}
	if err == nil && len(ids) > 0 {
		memberIds = ids
	} else {
		err = s.groupMapper.GetAllMemberId(msg.TargetId, &memberIds)
		if err != nil {
			fmt.Println("获取群成员id失败", err.Error())
			return err
		}
		//缓存构建
		err = redis.SetGroupMemberIds(msg.TargetId, memberIds)
		if err != nil {
			fmt.Println("缓存群成员id失败", err.Error())
		}
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
			_ = ws.WriteMessage(1, data)
		}()

	}
	return nil
}
