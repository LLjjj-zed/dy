package middleware

import (
	"testing"
)

func TestGetSnapshot(t *testing.T) {
	videopath := "bear.mp4"
	imagepath := "bear.png"
	err := GetSnapshotCmd(videopath, imagepath)
	if err != nil {
		t.Error(err)
	}
	//path := "C:\\Users\\violet\\Desktop\\bytedance\\douyin-demo\\public\\"
	//cmds := []string{
	//	"-i",
	//	path+"bear.mp4",
	//	"-ss 0:0:1",
	//	"-vframes",
	//	"1",
	//	path+"baer.png",
	//}
	//command := exec.Command("ffmpeg", cmds...)
	//err := command.Start()
	//if err != nil {
	//	t.Error(err)
	//}
}
