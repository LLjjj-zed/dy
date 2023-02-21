package Model

import (
	"douyin.core/config"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Video struct {
	Author        *User     `gorm:"-" json:"author"` // 视频作者信息
	UserID        int64     `json:"user_id"`         //用户id
	CommentCount  int64     `json:"comment_count"`   // 视频的评论总数
	CoverURL      string    `json:"cover_url"`       // 视频封面地址
	FavoriteCount int64     `json:"favorite_count"`  // 视频的点赞总数
	ID            int64     `json:"id"`              // 视频唯一标识
	IsFavorite    bool      `json:"is_favorite"`     // true-已点赞，false-未点赞
	PlayURL       string    `json:"play_url"`        // 视频播放地址
	Title         string    `json:"title"`           // 视频标题
	UserVideocode int64     `json:"videocode"`       //用户视频编号
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

// VideoDao 视频表数据操作结构体
type VideoDao struct {
}

// NewVideoDao 视频表数据操作结构体构造函数
func NewVideoDao() *VideoDao {
	return &VideoDao{}
}

// QueryVideoListLogin 查询用户视频关系表排除用户看过的视频
func (v *VideoDao) QueryVideoListLogin(last_time time.Time) ([]*Video, error) {
	var videos []*Video
	videos = make([]*Video, 0, config.MaxVideoList)
	err := DB.Model(&Video{}).Where("created_at<?", last_time).
		Order("created_at ASC").Limit(config.MaxVideoList).
		Select([]string{"id", "user_id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title"}).
		Find(&videos).Error

	if err != nil {
		return nil, err
	}
	return videos, nil
}

// QueryVideoListUnLogin 在未登录的状态下推送视频查询
func (v *VideoDao) QueryVideoListUnLogin(last_time time.Time) ([]*Video, error) {
	var videos []*Video
	videos = make([]*Video, 0, config.MaxVideoList)
	err := DB.Model(&Video{}).Where("created_at<?", last_time).
		Order("created_at ASC").Limit(config.MaxVideoList).
		Select([]string{"id", "user_id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title"}).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (v *VideoDao) AddAuthorInfoToFeedList(userid int64, videos *[]*Video) error {
	n := len(*videos)
	if videos == nil || n == 0 {
		return errors.New("不能使用空指针，视频列表为空")
	}
	userdao := NewUserInfoDao()
	//todo user upadteat
	//lasttime :=(*videos)[n-1].CreatedAt
	for i := 0; i < n; i++ {
		user, err := userdao.GetUserByuserID((*videos)[i].UserID)
		if err != nil {
			return err
		}
		//todo 获取是否点赞
		if userid != 0 {

		}
		(*videos)[i].Author = user
	}
	return nil
}

// PersistNewVideo 将视频数据持久化到数据库
func (v *VideoDao) PersistNewVideo(title string, userid int64, code int64, videoname, imagename string) error {
	playurl := GetUrl(videoname)
	coverurl := GetUrl(imagename)
	video := &Video{
		UserID:        userid,
		CommentCount:  0,
		CoverURL:      coverurl,
		FavoriteCount: 0,
		IsFavorite:    false,
		PlayURL:       playurl,
		Title:         title,
		UserVideocode: code,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err := DB.Table("videos").Create(video).Error
	if err != nil {
		return err
	}
	userdao := NewUserInfoDao()
	err = userdao.AddWorkCount(userid)
	if err != nil {
		return err
	}
	return nil
}

// GetUrl 获取url
func GetUrl(name string) string {
	var build strings.Builder
	build.WriteString("http://23.94.57.209:9000/")
	build.WriteString(name)
	url := build.String()
	return url
}

// GetUserVideoCode 获取用户视频序号吗，场景：用于用户将视频上传的时候生成视频文件名
func (v *VideoDao) GetUserVideoCode(userid int64) (int64, error) {
	var video Video
	result := DB.Select("user_videocode").Where("user_id=?", userid).Last(&video)
	if err := result.Error; err != nil {
		fmt.Println(err.Error())
		if err.Error() != "record not found" {
			return -1, err
		}
		return 0, nil
	}
	return video.UserVideocode + 1, nil
}

// QueryUserPublishList 查询用户发布列表
func (v *VideoDao) QueryUserPublishList(userid int64) (*VideoList, error) {
	var videos []*Video
	err := DB.Where("user_id=?", userid).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return &VideoList{Videos: videos}, nil
}

func QueryVideoById(vid int64) (Video, error) {
	var video Video
	err := DB.Where("id = ?", vid).First(&video).Error
	return video, err
}
