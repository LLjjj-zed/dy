package Model

import "gorm.io/gorm"

// User 用户信息表
type User struct {
	FollowCount     int64  `gorm:"column:follow_count" json:"follow_count"`     // 关注总数
	FollowerCount   int64  `gorm:"column:follower_count" json:"follower_count"` // 粉丝总数
	ID              int64  `gorm:"column:user_id" json:"user_id"`               // 用户id
	IsFollow        bool   `gorm:"column:is_follow" json:"is_follow"`           // true-已关注，false-未关注
	Name            string `gorm:"column:user_name" json:"user_name"`           // 用户名称
	Avatar          string `json:"avatar"`                                      // 用户头像
	BackgroundImage string `json:"background_image"`                            // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`                              // 喜欢数
	Signature       string `json:"signature"`                                   // 个人简介
	TotalFavorited  string `json:"total_favorited"`                             // 获赞数量
	WorkCount       int64  `json:"work_count"`                                  // 作品数
}

// UserInfoDao 用户信息数据操作结构体
type UserInfoDao struct {
}

// NewUserInfoDao 用户信息数据操作结构体构造函数
func NewUserInfoDao() *UserInfoDao {
	return &UserInfoDao{}
}

// GetUserByUserName 通过用户名查找用户
func (u *UserInfoDao) GetUserByUserName(username string) (*User, error) {
	var User User
	err := DB.Where("user_name= ? ", username).First(&User).Error
	if err != nil {
		return nil, err
	}
	return &User, nil
}

// InsertToUserInfoTable 将用户信息持久化到数据库
func (u *UserInfoDao) InsertToUserInfoTable(userid int64, username string) error {
	user := User{
		FollowCount:     0,
		FollowerCount:   0,
		ID:              userid,
		IsFollow:        false,
		Name:            username,
		Avatar:          "http://23.94.57.209:9000/browser/douyin/avatar.jpg",
		BackgroundImage: "http://23.94.57.209:9000/browser/douyin/avatar.jpg",
		FavoriteCount:   0,
		Signature:       "",
		TotalFavorited:  "0",
		WorkCount:       0,
	}
	return DB.Create(&user).Error
}

func (u *UserInfoDao) AddWorkCount(userid int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE users SET work_count=work_count+1 WHERE user_id = ?", userid).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetUserByuserID 通过用户ID查找用户
func (u *UserInfoDao) GetUserByuserID(userid interface{}) (*User, error) {
	var user User
	err := DB.Where("user_id=?", userid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserNameByUserID 通过用户id查询用户名，场景：用户上传视频的时候用于生成视频的文件名
func (u *UserInfoDao) GetUserNameByUserID(userid int64) (string, error) {
	var user User
	err := DB.Select("user_name").Where("user_id=?", userid).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Name, nil
}
