package dal

import (
	"douyin.core/config"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var dsn string

func ReadMysqlinfo() {
	info := config.Reader.GetStringMapString("mysql")
	dsn = info[config.Sqldsn]
}

var MaxVideoList = 15
var MaxLikeList = 15

func InitMysql() *gorm.DB {
	ReadMysqlinfo()
	fmt.Println(dsn)
	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		PrepareStmt:            true, //缓存预编译命令
		SkipDefaultTransaction: true, //禁用默认事务操作
		Logger:                 logger.Default.LogMode(logger.Info)})
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("Open MySQL failed")
	} else {
		logrus.Info("Open MySql successfully")
	}
	return DB
}
