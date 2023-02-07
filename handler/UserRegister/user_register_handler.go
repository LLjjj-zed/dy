package UserRegister

import (
	"douyin.core/Model"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"errors"
)

const (
	MaxUsernameLength = 20
	MaxPasswordLength = 18
	MinPasswordLength = 5
)



type UserRegisterReponse struct {
	Model.CommonResponse
	Token      string `json:"token"`      // 用户鉴权token
	UserID     int64  `json:"user_id"`    // 用户id
}

type PostUserLogin struct {
	username string
	password string
	userid   int64
	token    string
}

func (u *PostUserLogin) NewPostUserLogin(username,password string) *PostUserLogin {
	//创建对象，用于注册
	return &PostUserLogin{username: username,password: password}
}

func (u *PostUserLogin) PersistData()  {
	
}

func (u *PostUserLogin) CheckPost() error {
	if u.username == "" {
		return errors.New("用户名不能为空")
	}
	if len(u.username) > MaxUsernameLength {
		return errors.New("用户名过长")
	}

	if len(u.password) > MaxPasswordLength {
		return errors.New("密码过长")
	}
	if len(u.password) < MinPasswordLength{
		return errors.New("密码不能少于5位")
	}
	return nil
}




func RegisterOK(c *app.RequestContext)  {
	c.JSON(http.StatusOK,UserRegisterReponse{
		CommonResponse : Model.CommonResponse{
			StatusCode: 0,
		},
	})
}

func RegisterErr(c *app.RequestContext)  {
	c.JSON(http.StatusOK,UserRegisterReponse{
		CommonResponse : Model.CommonResponse{
			StatusCode: 1,
			StatusMsg: "fail",
		},
	})
}