package Model

import (
	"errors"
	"gorm.io/gorm"
)

// Type Like stored in gorm, contains video_id and user_id from Package Video ad User
type Like struct {
	gorm.Model

	VideoID int64 `json:"video_id"`
	Video   Video `gorm:"foreignKey:VideoID"`

	UserID int64 `json:"user_id"`
	User   User  `gorm:"foreignKey:UserID"`
}

type LikeDAO struct{}

func NewLikeDAO() *LikeDAO {
	return &LikeDAO{}
}

func (d LikeDAO) AddLike(userid int64, videoid int64) error {
	//点赞是否存在
	exist := DB.Where("user_id = ? AND video_id = ?", userid, videoid).First(&Like{})
	if exist != nil {
		return errors.New("Already liked")
	}
	var like = Like{
		VideoID: videoid,
		UserID:  userid,
	}
	return DB.Create(&like).Error
}

func (d LikeDAO) CancelLike(videoid int64, userid int64) error {
	var like Like
	return DB.Where(&Like{VideoID: videoid, UserID: userid}).Delete(&like).Error
}

func (d LikeDAO) QueryLikeList(userid int64) ([]Video, error) {
	var videoid []int64
	var videoList []Video
	err := DB.Where("user_id = ?", userid).Find(&videoid).Error
	for _, s := range videoid {
		video, _ := QueryVideoById(s)
		_ = append(videoList, video)
	}
	return videoList, err
}
