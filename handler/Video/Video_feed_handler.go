package Video

import (
	"douyin.core/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// 视频流回复结构体
type VideoFeedResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	NextTime   int64
	VideoList  []*Video
}

// 视频流处理函数，用于处理http请求
func VideoFeedHandler(c *gin.Context) {
	incoming_time := c.Query("latest_time")
	var last_time time.Time
	if incoming_time != "0" {
		//将字符串转换为数字,参数二为进制，参数三为位大小
		parseInt, _ := strconv.ParseInt(incoming_time, 10, 64)
		//转换为时间戳
		last_time = time.Unix(parseInt, 0)
	} else {
		last_time = time.Now()
	}
	//获取token并检查token是否存在
	token, exist := c.GetQuery("token")
	if exist {
		//token存在，向登录用户推送视频流
		LoginHandler(c, token, last_time)
	} else {
		//token不存在，向未登录用户推送视频流
		UnLoginHandeler(c, last_time)
	}
}

// 在用户未登录的状态下向用户推送视频的处理函数
func UnLoginHandeler(c *gin.Context, last_time time.Time) {
	dao := NewVideoDao()
	videolist, err := dao.QueryVideoListUnLogin(last_time)
	if err != nil {
		VideoFeedErr(c, err.Error())
		return
	}
	VideoFeedOK(c, videolist.Videos, last_time.Unix())
}

// 在用户已登录的状态下向用户推送视频的处理函数
func LoginHandler(c *gin.Context, token string, last_time time.Time) {
	claims, err := middleware.JwtParseUser(token)
	if err != nil {
		VideoFeedErr(c, err.Error())
		return
	}
	user_id := claims.Userid
	dao := NewVideoDao()
	videoList, err := dao.QueryVideoListLogin(user_id, last_time)
	if err != nil {
		VideoFeedErr(c, err.Error())
	}
	VideoFeedOK(c, videoList.Videos, last_time.Unix())
}

// 获取视频列表用于返回feed
func GetVideoList(c *gin.Context) *VideoList {
	var videos VideoList
	return &videos
}

// 返回正确信息
func VideoFeedOK(c *gin.Context, videos []*Video, next_time int64) {
	c.JSON(http.StatusOK, VideoFeedResponse{
		StatusCode: 0,
		StatusMsg:  "succese",
		NextTime:   next_time,
		VideoList:  videos,
	})
}

// 返回错误信息
func VideoFeedErr(c *gin.Context, ErrMasseage string) {
	c.JSON(http.StatusOK, VideoFeedResponse{
		StatusCode: 1,
		StatusMsg:  ErrMasseage,
	})
}
