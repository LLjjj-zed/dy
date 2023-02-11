package Video

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type VideoFeedResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	NextTime   time.Time
	VideoList  []*Video
}

func VideoFeedHandler(c *gin.Context) {

}

func UnLoginHandeler(c *gin.Context) {

}

func LoginHandler(c *gin.Context) {

}

func VideoFeedOK(c *gin.Context, videos []*Video) {
	c.JSON(http.StatusOK, VideoFeedResponse{
		StatusCode: 0,
		StatusMsg:  "succese",
		NextTime:   time.Now(),
		VideoList:  videos,
	})
}

func VideoFeedErr(c *gin.Context, ErrMasseage string) {
	c.JSON(http.StatusOK, VideoFeedResponse{
		StatusCode: 1,
		StatusMsg:  ErrMasseage,
	})
}
