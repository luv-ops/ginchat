package main

import (
	"GinChat/router"
	"GinChat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySql()
	utils.InitRedis()
	r := router.Router()
	r.Static("/static", "./static")
	r.Run(":8080")
}
