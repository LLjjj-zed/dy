package Model

import (
	user "douyin.core/handler/User"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	VideoId uint      `json:"videoId" gorm:"primary key;auto increment"` //视频ID
	Title   string    `json:"title" gorm:""`                             //视频标题
	UserID  user.User `json:"userID" gorm:"foreignKey:ID"`               //上传用户ID，外键关联至User
}
