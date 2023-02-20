package Like

import (
	"douyin.core/Model"
	"douyin.core/middleware"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetLikeList(c *gin.Context) {
	token, ok := c.GetQuery("token")
	if !ok {
		LikeListResponse(c, 1, "未能成功获取token，请重试", nil)
		return
	}
	_, err := middleware.JwtParseUser(token)
	if err != nil {
		LikeListResponse(c, 1, "token已过期，请重新登录", nil)
		return
	}
	//获取用户id
	userid, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if !ok {
		LikeListResponse(c, 1, "未能成功获取用户id，请重试", nil)
		return
	}
	dao := Model.LikeDAO{}
	likeList, err := dao.QueryLikeList(userid)
	if err != nil {
		LikeListResponse(c, 1, "获取点赞列表失败", nil)
		return
	}
	LikeListResponse(c, 0, "获取点赞列表成功", likeList)
}

func LikeListResponse(c *gin.Context, statuscode int64, statusmsg string, likeList []Model.Video) {
	c.JSON(200, gin.H{
		"status_code": statuscode,
		"status_msg":  statusmsg,
		"video_list":  likeList,
	})
}
