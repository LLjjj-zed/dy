package interaction

import (
	"douyin.core/Model"
	"douyin.core/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func FavoriteAction(c *gin.Context) {
	actionType := c.Query("action_type")
	token := c.Query("token")
	//获取用户信息
	usrid, err := middleware.JwtParseUser(token)
	if err != nil {
		return
	}
	videoid, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)

	if actionType == "1" {
		err := Model.AddLike(videoid, usrid.Userid)
		if err != nil {
			return
		}
	} else if actionType == "2" {
		err := Model.CancelLike(videoid, usrid.Userid)
		if err != nil {

		}
	} else {

	}
}

func GetFavList(c *gin.Context) {
	token := c.Query("token")
	_, err := middleware.JwtParseUser(token)
	if err != nil {
		return
	}
	usrid, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {

	}
	Videos, err := Model.QueryFavoriteList(usrid)
	c.JSON(http.StatusOK,
		gin.H{
			"StatusCode": 0,
			"StatusMsg":  "",
			"VideoList":  Videos,
		},
	)
}
