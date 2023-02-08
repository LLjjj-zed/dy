package user

import (
	"douyin.core/Model"
)

type UserInfoTable struct {
	FollowCount   int64  `gorm:"follow_count"`
	FollowerCount int64  `gorm:"follower_count"`
	ID            int64  `gorm:"id"`
	IsFollow      bool   `gorm:"is_follow"`
	Name          string `gorm:"name"`
}

type UserInfoDao struct {
}

func NewUserInfoDao() *UserInfoDao {
	return &UserInfoDao{}
}

func (u *UserInfoDao) GetUserByUserName(username string) UserInfoTable {
	var User UserInfoTable
	Model.DB.Where("username=?", username).First(&User)
	return User
}

func (u *UserInfoDao) InsertToUserInfoTable(userid int64, username string) error {
	user := UserInfoTable{
		FollowCount:   0,
		FollowerCount: 0,
		ID:            userid,
		IsFollow:      false,
		Name:          username,
	}
	return Model.DB.Create(&user).Error
}
