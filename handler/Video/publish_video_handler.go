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

type PublishVideoResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

func PublishVedioHandler(c *gin.Context) {
	title := c.PostForm("title")
	token, ok := c.GetQuery("token")
	if !ok {
		PublishvideoErr(c, "未能成功获取token，请重试")
		return
	}
	userclaim, err := middleware.JwtParseUser(token)
	if err != nil {
		PublishvideoErr(c, err.Error())
	}
	userid := userclaim.Userid
	file, err := c.FormFile("data")
	if err != nil {
		PublishvideoErr(c, err.Error())
		return
	}
	ext := filepath.Ext(file.Filename)
	if ext != ".mp4" {
		PublishvideoErr(c, "上传视频格式错误，请重试")
		return
	}
	videoDao := NewVideoDao()
	codeint, err := videoDao.GetUserVideoCode(userid)
	if err != nil {
		PublishvideoErr(c, err.Error())
		return
	}
	code := strconv.Itoa(int(codeint))
	var userinfo user.UserInfoDao
	username, err := userinfo.GetUserNameByUserID(userid)
	name := username
	if err != nil {
		PublishvideoErr(c, err.Error())
	}
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
	err = videoDao.PersistNewVideo(title, userid, &userinfo)
	if err != nil {
		PublishvideoErr(c, err.Error())
	}
	PublishVideoOk(c)
}

func PublishVideoOk(c *gin.Context) {
	c.JSON(http.StatusOK, &PublishVideoResponse{
		StatusCode: 0,
	})
}

func PublishvideoErr(c *gin.Context, errmeassage string) {
	c.JSON(http.StatusOK, &PublishVideoResponse{
		StatusCode: 1,
		StatusMsg:  errmeassage,
	})
}
