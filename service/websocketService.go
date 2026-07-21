package service

import (
	"GinChat/mapper"
	"GinChat/metrics"
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
type WsClient struct {
	Conn   *websocket.Conn
	Send   chan []byte //发送消息的信道
	UserId uint
	Closed chan struct{} //判断连接是否已经关闭的信道
	Done   chan struct{} //由于写协程阻塞，所以用于读协程通知
}

func (c *WsClient) writePump(exitFunc func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("WritePump panic: %v\n", r)
		}
		close(c.Send)
		close(c.Closed)
		exitFunc() //通知主协程
		_ = c.Conn.Close()
	}()
	for {
		select {
		case <-c.Done:
			return
		case msg, ok := <-c.Send:
			if !ok {
				return // Send 被关闭（自己关闭的，正常退出）
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				metrics.MsgSendFailTotal.Inc()
				fmt.Printf("写入失败，关闭连接: %v\n", err)
				return
			}
			metrics.MsgSendTotal.Inc()
		}

	}
}
func (c *WsClient) readPump(exitFunc func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("readLoop panic: %v\n", r)
		}
		close(c.Done) //通知写协程退出
		exitFunc()    //通知主协程

	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Printf("用户 %d 断开连接: %v\n", c.UserId, err)
			break
		}

		// 可以处理心跳消息
		msg := string(message)
		if msg == "ping" || msg == "heartbeat" {
			// 通过通道发送，由 WritePump 安全写出
			select {
			case c.Send <- []byte("pong"):
			default:
				fmt.Println("发送通道已满")
			}
		}

	}
}

var (
	wsConnMap = make(map[uint]*WsClient) //后期可以用redis 存储
	wsLock    sync.RWMutex               //读写锁,读可以共享，写独占
)

// WsConnectionHandler
// 管理用户map连接
func (s *WebsocketService) WsConnectionHandler(connect *websocket.Conn, userId uint) {
	client := &WsClient{
		Conn:   connect,
		Send:   make(chan []byte, 64),
		UserId: userId,
		Closed: make(chan struct{}),
		Done:   make(chan struct{}),
	}
	//加入map在线表
	wsLock.Lock()
	wsConnMap[userId] = client
	wsLock.Unlock()
	fmt.Println("用户上线：", userId)
	metrics.OnlineWSConn.Inc()
	_ = redis.SetUserOnline(userId, 1)

	exitChan := make(chan struct{}) //主协程使用的退出信号
	var exitOnce sync.Once
	exitFunc := func() { //只执行一次
		exitOnce.Do(func() {
			close(exitChan)
		})
	}
	//启动读协程
	go client.readPump(exitFunc)
	//启动写协程
	go client.writePump(exitFunc)

	// 等待任一协程退出
	<-exitChan
	//清理资源,断开连接

	//删除用户连接
	wsLock.Lock()
	delete(wsConnMap, userId)
	wsLock.Unlock()
	metrics.OnlineWSConn.Dec()
	_ = redis.SetUserOnline(userId, 0)
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
	select {
	case <-ws.Closed:
		return errors.New("用户已断开连接")
	case ws.Send <- data:
	default:
		fmt.Println("发送通道已满")
		return errors.New("发送通道已满")
	}
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
		data, _ := json.Marshal(msg)
		select {
		case <-ws.Closed:
			fmt.Println("发送用户已断开连接")
			continue
		case ws.Send <- data:
		default:
			fmt.Println("群消息丢弃")
			continue
		}

	}
	return nil
}
