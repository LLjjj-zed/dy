package middleware

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestInitminio(t *testing.T) {
	var ctx *gin.Context
	client := Initminio()
	videoname := "douyinbear.mp4"
	videopath := "../public/douyinbear.mp4"
	UploadVideoToMinio(ctx, client, videoname, videopath, "video")
}
