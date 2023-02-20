package middleware

//#include <stdlib.h>
//int startCmd(const char* cmd){
//	  return system(cmd);
//}
import "C"

import (
	"errors"
	"strings"
	"time"
	"unsafe"
)

var path = "C:\\Users\\violet\\Desktop\\bytedance\\douyin-demo\\public\\"

func GetSnapshotCmd(videoname, imagename string) error {
	var build strings.Builder
	build.WriteString("ffmpeg -i ")
	build.WriteString(path)
	build.WriteString(videoname)
	build.WriteString(" -ss 0:0:1 -vframes 1 ")
	build.WriteString(path)
	build.WriteString(imagename)
	cmd := build.String()
	cCmd := C.CString(cmd)
	defer C.free(unsafe.Pointer(cCmd))
	after := time.After(time.Second * 2)
	select {
	case <-after:
		return errors.New("timeout")
	default:
		status := C.startCmd(cCmd)
		if status != 0 {
			return errors.New("视频切截图失败")
		}
		return nil
	}
}
