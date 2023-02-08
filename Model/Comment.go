package Model

import (
	user "douyin.core/handler/User"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	CommentID   uint      `json:"commentID" gorm:"primary key;auto increment"`
	VideoId     Video     `json:"videoId" gorm:""`
	UserId      user.User `json:"userId" gorm:"foreignKey:ID"`
	CommentText string    `json:"commentText" gorm:"varchar(40);not null"`
}
