package Model

import (
	user "douyin.core/handler/UserInfo"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Video Video     `gorm:"foreignKey:VideoID"`
	User  user.User `gorm:"foreignKey:ID"`

	CommentID  uint   `json:"commentID" gorm:"primary key;auto increment"`
	VideoID    int64  `json:"videoId" `
	UsrID      int64  `json:"userId" `
	Content    string `json:"content" gorm:"varchar(40);not null"`
	CreateDate string `json:"create_date"`
}

func NewComment(Cmt Comment) error {
	err := DB.Create(Cmt).Error
	return err
}

func DeleteComment(commentID string) error {
	var comment Comment
	err := DB.Where("CommentID = ?", commentID).Delete(&comment).Error
	return err
}

func QueryCommentList(VideoID int64) ([]Comment, error) {
	var comments []Comment
	err := DB.Where("CommentId = ?", VideoID).Find(&comments).Error
	return comments, err
}
