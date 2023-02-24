package Model

import (
	"douyin.core/config"
	"douyin.core/dal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	DB = dal.DB
	config.Init()
	dal.InitMysql()
	dal.DB.AutoMigrate(User{}, UserLoginTable{}, Comment{}, Video{})
	dal.InitRedis()
}

func InitDB_test() {
	var err error
	DB, err = gorm.Open(mysql.Open(config.DBConnectString()), &gorm.Config{
		PrepareStmt:            true, //缓存预编译命令
		SkipDefaultTransaction: true, //禁用默认事务操作
		//Logger:                 logger.Default.LogMode(logger.Info), //打印sql语句
	})
	if err != nil {
		log.Fatal(err)
	}
	err = DB.AutoMigrate(&User{}, &Video{}, &Comment{}, &UserLoginTable{})
	if err != nil {
		log.Fatal(err)
	}
	dal.InitRedis()
}
