package UserRegister

import "douyin.core/Model"

// 用户登录表
type UserLoginTable struct {
	Id       int64  `grom:"primary_key"`
	UserId   int64  `grom:"notnull"`
	Username string `grom:"notnull"`
	Password string `grom:"notnull"`
}

type UserRigestDao struct {
}

func NewUserRigisterDao() *UserRigestDao {
	return &UserRigestDao{}
}

func (u *UserRigestDao) RegistUsertoDb(userid int64, username, password string) error {
	user := UserLoginTable{
		UserId:   userid,
		Username: username,
		Password: password,
	}
	return Model.DB.Create(&user).Error
}
