package Video

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 视频流回复结构体
type VideoFeedResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	NextTime   time.Time
	VideoList  []*Video
}

// 视频流处理函数，用于处理http请求
func VideoFeedHandler(c *gin.Context) {

}

// 在用户未登录的状态下向用户推送视频的处理函数
func UnLoginHandeler(c *gin.Context) {

}

// 在用户已登录的状态下向用户推送视频的处理函数
func LoginHandler(c *gin.Context) {

}

// 返回正确信息
func VideoFeedOK(c *gin.Context, videos []*Video) {
	c.JSON(http.StatusOK, VideoFeedResponse{
		StatusCode: 0,
		StatusMsg:  "succese",
		NextTime:   time.Now(),
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
