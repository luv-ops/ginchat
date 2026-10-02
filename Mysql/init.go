package Mysql

import (
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitMySql() {

	db, err := gorm.Open(mysql.Open(viper.GetString("Mysql.dsn")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Println("Mysql Autowired success")
	DB = db
}
