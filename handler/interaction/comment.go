package interaction

import (
	"douyin.core/Model"
	user "douyin.core/handler/User"
	"github.com/gin-gonic/gin"
)

// CommentResponse 评论回应结构
type CommentResponse struct {
	Comment
	Model.CommonResponse // 返回状态描述
}

// Users 评论用户信息
type Users = user.User

// Comment 返回评论信息
type Comment struct {
	Content    string `json:"content"`     // 评论内容
	CreateDate string `json:"create_date"` // 评论发布日期，格式 mm-dd
	ID         int64  `json:"id"`          // 评论id
	User       Users  `json:"user"`        // 评论用户信息
}

// CommentRequest 评论请求结构
type CommentRequest struct {
	Token       string `json:"token,omitempty"`
	VideoId     string `json:"video_id,omitempty"`
	ActionType  string `json:"action_type,omitempty"`
	CommentText string `json:"comment_text,omitempty"`
	CommentID   string `json:"comment_id,omitempty"`
}

func CommentAction(c *gin.Context) {
	var Req CommentRequest
	_ = c.ShouldBindJSON(&Req)
	//todo: Validate Token

	if Req.ActionType == "1" {
		err := Model.NewComment(&Req)
		if err != nil {
			return
		}
	} else if Req.ActionType == "2" {
		err := Model.DeleteComment(Req.CommentID)
		if err != nil {
			//
		}
	}
}

func DeleteComment(c *gin.Context) {}
