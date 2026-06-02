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

func Router() *gin.Engine {
	r := gin.Default()
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
			public.POST("/register", controller.Register)
			public.POST("/login", controller.Login)
		}
		//需要token
		auth := userGroup.Group("/")
		auth.Use(middleware.JwtAuth())
		{
			auth.GET("/list", controller.GetUserList)
			auth.GET("/info", controller.UserInfo)
			auth.POST("/delete", controller.DeleteUser)
			auth.POST("/update", controller.UpdateUser)

		}
	}
	//好友模块
	friendGroup := r.Group("/friend")
	{
		friendGroup.Use(middleware.JwtAuth())
		friendGroup.POST("/add", controller.AddFriend)
		requestGroup := friendGroup.Group("/requests")
		{
			requestGroup.GET("", controller.RequestList)
			requestGroup.GET("/unread", controller.UnReadCount)
			requestGroup.POST("/hasRead", controller.HasRead)
		}

		friendGroup.POST("/accept/:friendId", controller.Accept)
		friendGroup.POST("/reject/:friendId", controller.Reject)
		friendGroup.GET("/list", controller.GetFriendList)

	}
	//websocket模块
	//websocket  ws://127.0.0.1:8080/socket通过这个路径进入websocket
	wsGroup := r.Group("/ws")
	{
		wsGroup.Use(middleware.JwtAuth(true))
		wsGroup.GET("/:token", controller.ConnectWs)
	}
	//聊天模块
	chatGroup := r.Group("/chat")
	{
		chatGroup.Use(middleware.JwtAuth())
		chatGroup.POST("/send", controller.Send)
	}
	//会话模块
	conversationGroup := r.Group("/conversation")
	{
		conversationGroup.Use(middleware.JwtAuth())
		conversationGroup.GET("/list", controller.ConversationList)
		conversationGroup.POST("/unreadClear/:peerId", controller.ClearUnreadCount)
	}
	//消息模块
	messageGroup := r.Group("/message")
	{
		messageGroup.Use(middleware.JwtAuth())
		messageGroup.GET("/list", controller.GetMessage)
	}
	//群组模块
	groupGroup := r.Group("/group")
	{
		groupGroup.Use(middleware.JwtAuth())
		groupGroup.POST("/create", controller.CreateGroup)
		groupGroup.POST("/invite", controller.InviteGroup)
		groupGroup.GET("/detail/:groupId", controller.GroupDetail)
		groupGroup.GET("/members/:groupId", controller.GroupMembers)
	}
	//上传模块
	uploadGroup := r.Group("/upload")
	{
		uploadGroup.Use(middleware.JwtAuth())
		uploadGroup.POST("/file", controller.UploadFile)
	}
	return r
}
