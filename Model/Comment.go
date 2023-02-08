package Model

import (
	user "douyin.core/handler/User"
	"douyin.core/handler/interaction"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	CommentID  uint      `json:"commentID" gorm:"primary key;auto increment"`
	VideoId    Video     `json:"videoId" gorm:""`
	UserId     user.User `json:"userId" gorm:"foreignKey:ID"`
	Content    string    `json:"content" gorm:"varchar(40);not null"`
	CreateDate string    `json:"create_date"`
}

func NewComment(Req *interaction.CommentRequest) error {
	return nil
}

func DeleteComment(commentID string) error {
	var comment Comment
	err := DB.Where("CommentID = ?", commentID).Delete(&comment).Error
	return err
}
