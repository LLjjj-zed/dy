package Video

import (
	"douyin.core/Model"
	user "douyin.core/handler/User"
)

// Video
type Video struct {
	Author        *user.User `json:"author"`         // 视频作者信息
	CommentCount  int64      `json:"comment_count"`  // 视频的评论总数
	CoverURL      string     `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64      `json:"favorite_count"` // 视频的点赞总数
	ID            int64      `json:"id"`             // 视频唯一标识
	IsFavorite    bool       `json:"is_favorite"`    // true-已点赞，false-未点赞
	PlayURL       string     `json:"play_url"`       // 视频播放地址
	Title         string     `json:"title"`          // 视频标题
	UserVideocode int64      `json:"videodode"`      //用户视频编号
}

type VideoDao struct {
}

func NewVideoDao() *VideoDao {
	return &VideoDao{}
}

func (v *VideoDao) QueryVideoby() {

}

func (v *VideoDao) PersistNewVideo(title string, userid int64, user *user.UserInfoDao) error {
	userinfo, err := user.GetUserByuserID(userid)
	if err != nil {
		return err
	}
	video := &Video{
		Author:        userinfo,
		CommentCount:  0,
		CoverURL:      "",
		FavoriteCount: 0,
		ID:            0,
		IsFavorite:    false,
		PlayURL:       "",
		Title:         title,
		UserVideocode: 0,
	}
	return Model.DB.Create(video).Error
}

func (v *VideoDao) GetUserVideoCode(userid int64) (int64, error) {
	var video Video
	err := Model.DB.Where("").Last(&video).Error
	if err != nil {
		return -1, err
	}
	if video.ID == 0 {
		return 0, nil
	}
	return video.UserVideocode + 1, nil
}
