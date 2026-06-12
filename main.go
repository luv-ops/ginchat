package main

import (
	"GinChat/Autowired"
	"GinChat/Mysql"
	"GinChat/config"
	"GinChat/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	Mysql.InitMySql()
	redis.InitRedis()
	//依赖注入
	route := Autowired.InitAll()
	r := gin.Default()
	route.Setup(r)
	r.Static("/static", "./static")
	r.Run(":8080")
}
