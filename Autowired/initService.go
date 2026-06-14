package Autowired

import (
	"GinChat/MQ"
	"GinChat/Mysql"
	"GinChat/service"
)

var (
	UserService         *service.UserService
	MessageService      *service.MessageService
	GroupService        *service.GroupService
	FriendService       *service.FriendService
	ConversationService *service.ConversationService
	ChatService         *service.ChatService
	WebsocketService    *service.WebsocketService
)

func InitService() {
	// ✅ 先初始化 WebsocketService（它实现 MessageSender 接口）
	WebsocketService = service.NewWebsocketService(GroupMapper, UserMapper)

	// ✅ 再初始化依赖 MessageSender 的 Service
	UserService = service.NewUserService(UserMapper)
	MessageService = service.NewMessageService(MessageMapper)
	GroupService = service.NewGroupService(GroupMapper, ConversationMapper, Mysql.DB)

	// ✅ WebsocketService 实现了 MessageSender 接口，可以直接传入
	FriendService = service.NewFriendService(FriendMapper, UserMapper, WebsocketService, Mysql.DB, MQ.GlobalKafkaCli)

	ConversationService = service.NewConversationService(ConversationMapper)
	ChatService = service.NewChatService(UserMapper, ConversationMapper, MessageMapper, WebsocketService, Mysql.DB, MQ.GlobalKafkaCli)
	//FriendService，ChatService实现了MessageHandler接口
	//注入到MQ
	MQ.FriReqHandler = FriendService
	MQ.MsgHandler = ChatService
}
