package Video

import (
	"douyin.core/Model"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// FeedResponse 视频流回复结构体
type FeedResponse struct {
	StatusCode int64          `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string         `json:"status_msg"`  // 返回状态描述
	NextTime   int64          `json:"next_time"`
	VideoList  []*Model.Video `json:"video_list"`
	latestTime time.Time
}

// VideoFeedHandler 视频流处理函数，用于处理http请求
func VideoFeedHandler(c *gin.Context) {
	incomingTime := c.Query("latest_time")
	var lastTime time.Time
	if incomingTime != "0" {
		//将字符串转换为数字,参数二为进制，参数三为位大小
		parseInt, _ := strconv.ParseInt(incomingTime, 10, 64)
		//转换为时间戳
		lastTime = time.Unix(0, parseInt*1e6)
	} else {
		lastTime = time.Now()
	}
	//获取token并检查token是否存在
	token, exist := c.GetQuery("token")
	//fmt.Println(token)
	if exist {
		//token存在，向登录用户推送视频流
		LoginHandler(c, token, lastTime)
	} else {
		//token不存在，向未登录用户推送视频流
		UnLoginHandler(c, lastTime)
	}
}

// UnLoginHandler 在用户未登录的状态下向用户推送视频的处理函数
func UnLoginHandler(c *gin.Context, lastTime time.Time) {
	dao := Model.NewVideoDao()
	videos, err := dao.QueryVideoListUnLogin(lastTime)
	if err != nil {
		FeedErr(c, err.Error())
		return
	}
	//判断视频列表长度，防止panic
	if len(videos) <= 0 {
		FeedErr(c, errors.New("当前无最新视频").Error())
		return
	}
	err = dao.AddAuthorInfoToFeedList(0, &videos)
	if err != nil {
		FeedErr(c, err.Error())
		return
	}
	FeedOK(c, videos, lastTime.Unix())
}

// LoginHandler 在用户已登录的状态下向用户推送视频的处理函数
func LoginHandler(c *gin.Context, token string, lastTime time.Time) {
	//claims, err := middleware.ParseToken(token)
	//if err != nil {
	//	FeedErr(c, errors.New("token not exist").Error())
	//	return
	//}
	//todo 推送user未观看过的视频
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	//fmt.Println(userId)
	dao := Model.NewVideoDao()
	videos, err := dao.QueryVideoListLogin(lastTime)
	if err != nil {
		FeedErr(c, err.Error())
		return
	}
	err = dao.AddAuthorInfoToFeedList(userId, &videos)
	if err != nil {
		FeedErr(c, err.Error())
		return
	}
	//fmt.Println(videos[0])
	//DemoFeed(c)
	FeedOK(c, videos, lastTime.Unix())
}

//func DemoFeed(c *gin.Context)  {
//	c.JSON(http.StatusOK, FeedResponse{
//		StatusCode: 0,
//		VideoList: Model.DemoVideo,
//		NextTime:  time.Now().Unix(),
//	})
//}

// FeedOK 返回正确信息
func FeedOK(c *gin.Context, videos []*Model.Video, nextTime int64) {
	c.JSON(http.StatusOK, FeedResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		NextTime:   nextTime,
		VideoList:  videos,
	})
}

// FeedErr 返回错误信息
func FeedErr(c *gin.Context, ErrMessage string) {
	c.JSON(201, FeedResponse{
		StatusCode: 1,
		StatusMsg:  ErrMessage,
	})
}
