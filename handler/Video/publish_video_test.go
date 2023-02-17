package Video

import (
	"fmt"
	"testing"
)

func TestGetFilename(t *testing.T) {
	videoname := GetFilename("lljjj", "0", ".mp4")
	imagename := GetFilename("lljjj", "0", ".jpg")
	fmt.Println(videoname)
	fmt.Println(imagename)

}
