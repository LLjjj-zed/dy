package utils

import (
	"fmt"
	"testing"
)

func TestInitminio(t *testing.T) {
	//var ctx *gin.Context
	client := Initminio()
	fmt.Println(client)
	videoname := "douyinbear.mp4"
	videopath := "../public/douyinbear.mp4"
	err := UploadVideoToMinio(client, videoname, videopath, "video")
	if err != nil {
		t.Error(err)
	}
}
