package main

import (
	"GinChat/Mysql"
	"GinChat/config"
	"GinChat/redis"
	"GinChat/router"
)

func main() {
	config.InitConfig()
	Mysql.InitMySql()
	redis.InitRedis()
	r := router.Router()
	r.Static("/static", "./static")
	r.Run(":8080")
}
