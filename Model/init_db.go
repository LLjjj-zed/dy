package Model

import (
	"douyin.core/config"

	"douyin.core/dal"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	DB = dal.DB
	config.Init()
	dal.InitMysql()
	dal.DB.AutoMigrate(User{}, UserLoginTable{}, Comment{}, Video{})
	dal.InitRedis()
}
