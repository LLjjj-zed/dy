package Comment

import (
	"douyin.core/Model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommentListResponse struct {
	CommentList []Model.Comment `json:"comment_list"` // 评论列表
	StatusCode  int64           `json:"status_code"`  // 状态码，0-成功，其他值-失败
	StatusMsg   string          `json:"status_msg"`   // 返回状态描述
}

func GetCommentList(c *gin.Context) {
	var err error
	dao := Model.CommentDao{}
	cmtList, err := dao.GetCommentList(c.Query("video_id"))
	if err != nil {
		CommentBadResponse(c, err.Error())
		return
	}
	GetListResponse(c, cmtList)

}

func GetListResponse(c *gin.Context, comments []Model.Comment) {
	c.JSON(http.StatusOK, CommentListResponse{
		CommentList: comments,
		StatusCode:  0,
		StatusMsg:   "succeess",
	})
}
