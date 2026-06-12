package router

import (
	"GinChat/controller"
	"GinChat/docs"
	"GinChat/middleware"
	"GinChat/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	userController         controller.IUserController
	chatController         controller.IChatController
	friendController       controller.IFriendController
	websocketController    controller.IWebsocketController
	conversationController controller.IConversationController
	messageController      controller.IMessageController
	groupController        controller.IGroupController
}

func NewRouter(uC controller.IUserController, chatC controller.IChatController,
	friendC controller.IFriendController, websocketC controller.IWebsocketController,
	conversationC controller.IConversationController, messageC controller.IMessageController,
	groupC controller.IGroupController) *Router {
	return &Router{
		userController:         uC,
		chatController:         chatC,
		friendController:       friendC,
		websocketController:    websocketC,
		conversationController: conversationC,
		messageController:      messageC,
		groupController:        groupC,
	}
}
func (R *Router) Setup(r *gin.Engine) {

	r.Use(middleware.Cors())
	r.GET("", func(c *gin.Context) {
		utils.Ok2(c, "成功部署")
	})
	//引入swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//用户模块
	userGroup := r.Group("/user")
	{
		//不需要token
		public := userGroup.Group("/")
		{
			public.POST("/register", R.userController.Register)
			public.POST("/login", R.userController.Login)
		}
		//需要token
		auth := userGroup.Group("/")
		auth.Use(middleware.JwtAuth())
		{
			auth.GET("/list", R.userController.GetUserList)
			auth.GET("/info", R.userController.UserInfo)
			auth.POST("/delete", R.userController.DeleteUser)
			auth.POST("/update", R.userController.UpdateUser)

		}
	}
	//好友模块
	friendGroup := r.Group("/friend")
	{
		friendGroup.Use(middleware.JwtAuth())
		friendGroup.POST("/add", R.friendController.AddFriend)
		requestGroup := friendGroup.Group("/requests")
		{
			requestGroup.GET("", R.friendController.RequestList)
			requestGroup.GET("/unread", R.friendController.UnReadCount)
			requestGroup.POST("/hasRead", R.friendController.HasRead)
		}

		friendGroup.POST("/accept/:friendId", R.friendController.Accept)
		friendGroup.POST("/reject/:friendId", R.friendController.Reject)
		friendGroup.GET("/list", R.friendController.GetFriendList)

	}
	//websocket模块
	//websocket  ws://127.0.0.1:8080/socket通过这个路径进入websocket
	wsGroup := r.Group("/ws")
	{
		wsGroup.Use(middleware.JwtAuth(true))
		wsGroup.GET("/:token", R.websocketController.ConnectWs)
	}
	//聊天模块
	chatGroup := r.Group("/chat")
	{
		chatGroup.Use(middleware.JwtAuth())
		chatGroup.POST("/send", R.chatController.Send)
	}
	//会话模块
	conversationGroup := r.Group("/conversation")
	{
		conversationGroup.Use(middleware.JwtAuth())
		conversationGroup.GET("/list", R.conversationController.ConversationList)
		conversationGroup.POST("/unreadClear/:peerId", R.conversationController.ClearUnreadCount)
	}
	//消息模块
	messageGroup := r.Group("/message")
	{
		messageGroup.Use(middleware.JwtAuth())
		messageGroup.GET("/list", R.messageController.GetMessage)
	}
	//群组模块
	groupGroup := r.Group("/group")
	{
		groupGroup.Use(middleware.JwtAuth())
		groupGroup.POST("/create", R.groupController.CreateGroup)
		groupGroup.POST("/invite", R.groupController.InviteGroup)
		groupGroup.GET("/detail/:groupId", R.groupController.GroupDetail)
		groupGroup.GET("/members/:groupId", R.groupController.GroupMembers)
	}
	//上传模块
	uploadGroup := r.Group("/upload")
	{
		uploadGroup.Use(middleware.JwtAuth())
		uploadGroup.POST("/file", controller.UploadFile)
	}
}
