package Autowired

import "GinChat/controller"

var (
	UserController         *controller.UserController
	MessageController      *controller.MessageController
	GroupController        *controller.GroupController
	FriendController       *controller.FriendController
	ConversationController *controller.ConversationController
	ChatController         *controller.ChatController
	WebsocketController    *controller.WebsocketController
)

func InitController() {
	UserController = controller.NewUserController(UserService)
	MessageController = controller.NewMessageController(MessageService)
	GroupController = controller.NewGroupController(GroupService)
	FriendController = controller.NewFriendController(FriendService)
	ConversationController = controller.NewConversationController(ConversationService)
	ChatController = controller.NewChatController(ChatService)
	WebsocketController = controller.NewWebsocketController(WebsocketService)
}
