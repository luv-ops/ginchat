package Autowired

import "GinChat/router"

func InitAll() *router.Router {
	InitMapper()
	InitService()
	InitController()
	return router.NewRouter(
		UserController,
		ChatController,
		FriendController,
		WebsocketController,
		ConversationController,
		MessageController,
		GroupController,
	)
}
