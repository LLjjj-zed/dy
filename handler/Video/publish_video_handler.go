package Video

import (
	user "douyin.core/handler/User"
	"douyin.core/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// 用户投稿回复结构体
type PublishVideoResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

// 用户投稿处理函数，用于处理http请求
func PublishVedioHandler(c *gin.Context) {
	//从请求中获取视频标题和token
	title := c.PostForm("title")
	token, ok := c.GetQuery("token")
	if !ok {
		PublishvideoErr(c, "未能成功获取token，请重试")
		return
	}
	//从token中解析出用户id
	userclaim, err := middleware.JwtParseUser(token)
	if err != nil {
		PublishvideoErr(c, err.Error())
	}
	userid := userclaim.Userid
	//从请求中获取时视频数据
	file, err := c.FormFile("data")
	if err != nil {
		PublishvideoErr(c, err.Error())
		return
	}
	//从视频数据中获取视频后缀格式
	ext := filepath.Ext(file.Filename)
	if ext != ".mp4" {
		PublishvideoErr(c, "上传视频格式错误，请重试")
		return
	}
	videoDao := NewVideoDao()
	//获取用户视频序号，用于生成视频文件名
	codeint, err := videoDao.GetUserVideoCode(userid)
	if err != nil {
		PublishvideoErr(c, err.Error())
		return
	}
	//将用户视频序号转换成字符串
	code := strconv.Itoa(int(codeint))
	var userinfo user.UserInfoDao
	//获取用户名，用于生成视频文件名
	username, err := userinfo.GetUserNameByUserID(userid)
	name := username
	if err != nil {
		PublishvideoErr(c, err.Error())
	}
	//将视频持久化到本地，使用strings.Builder替换+提高性能
	var build strings.Builder
	build.WriteString(name)
	build.WriteString("code")
	build.WriteString(code)
	build.WriteString(ext)
	filename := build.String()
	path := filepath.Join("../puclic", filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		PublishvideoErr(c, err.Error())
		return
	}
	//将视频信息持久化到数据库
	err = videoDao.PersistNewVideo(title, userid, &userinfo)
	if err != nil {
		PublishvideoErr(c, err.Error())
	}
	PublishVideoOk(c)
}

// 返回正确信息
func PublishVideoOk(c *gin.Context) {
	c.JSON(http.StatusOK, &PublishVideoResponse{
		StatusCode: 0,
	})
}

// 返回错误信息
func PublishvideoErr(c *gin.Context, errmeassage string) {
	c.JSON(http.StatusOK, &PublishVideoResponse{
		StatusCode: 1,
		StatusMsg:  errmeassage,
	})
}
