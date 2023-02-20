package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"strings"
)

type Mysql struct {
	Host      string
	Port      int
	Database  string
	Username  string
	Password  string
	Charset   string
	ParseTime bool `toml:"parse_time"`
	Loc       string
}

type Server struct {
	IP   string
	Port int
}

type Config struct {
	DB     Mysql `toml:"mysql"`
	Server `toml:"server"`
}

var Secret = "tiktok"

var Info Config

// 包初始化加载时候会调用的函数
func init() {
	if _, err := toml.DecodeFile("C:\\Users\\violet\\Desktop\\douyin-demo\\config\\config.toml", &Info); err != nil {
		log.Fatal(err)
	}
	//去除左右的空格
	strings.Trim(Info.Server.IP, " ")
	strings.Trim(Info.DB.Host, " ")
}

var MaxVideoList = 15
var MaxLikeList = 15

// DBConnectString 填充得到数据库连接字符串
func DBConnectString() string {
	arg := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		Info.DB.Username, Info.DB.Password, Info.DB.Host, Info.DB.Port, Info.DB.Database,
		Info.DB.Charset, Info.DB.ParseTime, Info.DB.Loc)
	log.Println(arg)
	return arg
}
