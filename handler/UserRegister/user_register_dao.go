package UserRegister

import (
	"douyin.core/Model"
	"errors"
)

// 用户登录表
type UserLoginTable struct {
	Id       int64  `grom:"primary_key"`
	UserId   int64  `grom:"notnull"`
	Username string `grom:"notnull"`
	Password string `grom:"notnull"`
}

type UserRigestDao struct {
}

func NewUserLoginTable(username, password string) *UserLoginTable {
	return &UserLoginTable{Username: username, Password: password}
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

func (u UserRigestDao) QueryUserLogin(username, password string, login *UserLoginTable) error {
	err := Model.DB.Where("username=?", username).First(&login).Error
	if err != nil {
		return err
	}
	if login.Password != password {
		err = errors.New("密码错误")
		return err
	}
	return nil
}
