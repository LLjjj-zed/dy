package middleware

import "testing"

func TestInitminio(t *testing.T) {
	client, ctx := Initminio()
	videoname := "douyinbear.mp4"
	videopath := "../public/douyinbear.mp4"
	UploadFileToMinio(client, videoname, videopath, ctx)
}
