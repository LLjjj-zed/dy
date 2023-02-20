package Comment

import (
	user "douyin.core/Model"
	"douyin.core/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PublishCommentResponse struct {
	StatusCode int64         `json:"status_code"`
	StatusMsg  string        `json:"status_msg"`
	Comment    *user.Comment `json:"comment"`
}

func CommentActionHandler(c *gin.Context) {
	token, ok := c.GetQuery("token")
	if !ok {
		//ERR
		return
	}
	var err error
	userclaim, err := middleware.JwtParseUser(token)
	if err != nil {
		CommentBadResponse(c, err.Error())
	}
	userid := userclaim.Userid
	actionType := c.PostForm("action_type")
	cmtDao := user.NewCommentDao()
	var userinfo user.UserInfoDao

	if actionType == "1" {
		content, ok := c.GetQuery("comment_text")
		if !ok {
			CommentBadResponse(c, "comment_text is required")
			return
		}
		err, newcmt := cmtDao.AddComment(userid, content, &userinfo)
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
