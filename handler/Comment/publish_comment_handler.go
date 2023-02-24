package Comment

import (
	user "douyin.core/Model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PublishCommentResponse struct {
	StatusCode int64         `json:"status_code"`
	StatusMsg  string        `json:"status_msg"`
	Comment    *user.Comment `json:"comment"`
}

func CommentActionHandler(c *gin.Context) {
	userid, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	actionType := c.Query("action_type")
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	cmtDao := user.NewCommentDao()
	var userinfo user.UserInfoDao

	if actionType == "1" {
		content, ok := c.GetQuery("comment_text")
		if !ok {
			CommentBadResponse(c, "comment_text is required")
			return
		}
		err, newcmt := cmtDao.AddComment(userid, content, &userinfo, videoId)
		if err != nil {
			CommentBadResponse(c, err.Error())
			return
		}
		CommentSuccessResponse(c, &newcmt)
	} else if actionType == "2" {
		cmtid, ok := c.GetQuery("comment_id")
		if !ok {
			CommentBadResponse(c, "comment_id is required")
			return
		}
		err := cmtDao.DeleteComment(cmtid)
		if err != nil {
			CommentBadResponse(c, err.Error())
			return
		}
	} else {
		CommentBadResponse(c, "Invalid action type")
		return
	}
}

func CommentBadResponse(c *gin.Context, errmsg string) {
	c.JSON(http.StatusOK, PublishCommentResponse{
		StatusCode: 1,
		StatusMsg:  errmsg,
	})
}

func CommentSuccessResponse(c *gin.Context, comment *user.Comment) {
	c.JSON(http.StatusOK, PublishCommentResponse{
		StatusCode: 0,
		StatusMsg:  "OK",
		Comment:    comment,
	})
}
