package Model

import (
	user "douyin.core/handler/User"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	User user.User `gorm:"foreignKey:ID; references:ID"`

	VideoId uint   `json:"videoId" gorm:"primary key;auto increment"` //视频ID
	Title   string `json:"title" gorm:""`                             //视频标题
	UserID  int64  `json:"userID"`                                    //上传用户ID，外键关联至User
}

func QueryVideoById(vid int64) (Video, error) {
	var video Video
	err := DB.Where("videoid = ?", vid).First(&video).Error
	return video, err
}
