package Model

//import (
//	"douyin.core/handler/User"
//)
//
//type Favorite struct {
//	Video Video     `gorm:"foreignKey:VideoId"`
//	User  User.User `gorm:"foreignKey:UserId"`
//
//	VideoId int64 `json:"video_id"`
//	UserId  int64 `json:"user_id"`
//}
//
//func AddLike(videoid int64, userid int64) error {
//	return DB.Create(Favorite{
//		VideoId: videoid,
//		UserId:  userid,
//	}).Error
//}
//
//func CancelLike(videoid int64, userid int64) error {
//	var like Favorite
//	return DB.Where(&Favorite{VideoId: videoid, UserId: userid}).Delete(&like).Error
//}
//
//func QueryFavoriteList(userid int64) ([]Video, error) {
//	var videoid []int64
//	var videoList []Video
//	err := DB.Where("UserId = ?", userid).Find(&videoid).Error
//	for _, s := range videoid {
//		video, _ := QueryVideoById(s)
//		_ = append(videoList, video)
//	}
//	return videoList, err
//}
