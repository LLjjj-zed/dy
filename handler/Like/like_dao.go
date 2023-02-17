package Like

import (
	"douyin.core/Model"
	"douyin.core/handler/Video"
	"gorm.io/gorm"
)

// Type Like stored in gorm, contains video_id and user_id from Package Video ad User
type Like struct {
	gorm.Model
	VideoID string
	UserID  string
}

type LikeDAO struct{}

func NewLikeDAO() *LikeDAO {
	return &LikeDAO{}
}

func (d LikeDAO) AddLike(userid int64, videoid int64) error {
	return Model.DB.Create(Like{
		VideoID: string(videoid),
		UserID:  string(userid),
	}).Error
}

func (d LikeDAO) CancelLike(videoid int64, userid int64) error {
	var like Like
	return Model.DB.Where(&Like{VideoID: string(videoid), UserID: string(userid)}).Delete(&like).Error
}

func (d LikeDAO) QueryLikeList(userid int64) ([]Video.Video, error) {
	var videoid []int64
	var videoList []Video.Video
	err := Model.DB.Where("user_id = ?", userid).Find(&videoid).Error
	for _, s := range videoid {
		video, _ := Video.QueryVideoById(s)
		_ = append(videoList, video)
	}
	return videoList, err
}
