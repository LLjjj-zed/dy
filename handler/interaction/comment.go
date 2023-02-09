package interaction

import (
	"douyin.core/Model"
	user "douyin.core/handler/UserInfo"
	"douyin.core/middleware"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
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

// CommentAction 新增/删除评论
func CommentAction(c *gin.Context) {
	//将请求映射到结构中
	var Req CommentRequest
	_ = c.ShouldBindJSON(&Req)
	//获取用户信息
	var err error = nil
	usrid, err := middleware.MwuserId(c)
	userDao := user.NewUserInfoDao()
	usr, err := userDao.GetUserByuserID(usrid)
	if err != nil {
		return
	}
	//检查操作类型，1为新增，2为删除
	if Req.ActionType == "1" {
		//生成存储结构
		VideoId, _ := strconv.Atoi(Req.VideoId)
		comment := Model.Comment{
			Content:    Req.CommentText,
			UsrID:      usr.ID,
			VideoID:    int64(VideoId),
			CreateDate: time.Now().String(),
		}
		err = Model.NewComment(comment)
		if err != nil {
			return
		}
		respCmt := &Comment{
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.String(),
			ID:         int64(comment.CommentID),
			User:       *usr,
		}
		CommentOK(c, respCmt)
	} else if Req.ActionType == "2" {
		err = Model.DeleteComment(Req.CommentID)
		if err != nil {
			//
		}
	} else {
		err = errors.New("invalid action type")
		return
	}
}

// GetCommentList 获取评论列表
func GetCommentList(c *gin.Context) {
	_, err := middleware.MwuserId(c)
	if err != nil {
		return
	}
	VideoID, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Model.CommonResponse{
			StatusCode: -1,
			StatusMsg:  "comment videoId json invalid",
		})
		log.Println("CommentController-Comment_List: return videoId json invalid") //视频id格式有误
		return
	}
	CmtList, err := Model.QueryCommentList(VideoID)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK,
		gin.H{
			"StatusCode":  0,
			"StatusMsg":   "",
			"CommentList": CmtList,
		},
	)
}

func CommentOK(c *gin.Context, cmt *Comment) {
	c.JSON(http.StatusOK, CommentResponse{
		Comment{
			Content:    cmt.Content,
			CreateDate: cmt.CreateDate,
			ID:         cmt.ID,
			User:       cmt.User,
		},
		Model.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
	})
}
