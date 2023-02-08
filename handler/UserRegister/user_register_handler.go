package UserRegister

import (
	"douyin.core/Model"
	"douyin.core/middleware"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	MaxUsernameLength = 20
	MaxPasswordLength = 18
	MinPasswordLength = 5
)

type UserRegisterReponse struct {
	Model.CommonResponse
	Token  string `json:"token"`   // 用户鉴权token
	UserID int64  `json:"user_id"` // 用户id
}

// 用户登录校验
type PostUserLogin struct {
	Username string
	Password string
	Userid   int64
	Token    string
}

func UserLoginHandler(c *gin.Context) {
	//获取用户名和密码
	username := c.Query("username")
	get, exists := c.Get("password")
	if !exists {
		RegisterErr(c, "密码不能为空")
		return
	}
	password := get.(string)
	//创建新对象
	newPostUserLogin := NewPostUserLogin(username, password)
	//注册新用户
	err := newPostUserLogin.Regist()
	if err != nil {
		RegisterErr(c, err.Error())
		return
	}
	RegisterOK(c, newPostUserLogin)
}

// 注册新用户
func (u *PostUserLogin) Regist() error {
	//校验参数
	err := u.CheckPost()
	if err != nil {
		return err
	}
	//持久化到数据库
	err = u.PersistData()
	if err != nil {
		return err
	}
	//获取token
	err = u.SetToken()
	if err != nil {
		return err
	}
	//获取userid

	return nil
}

// 创建对象，用于注册
func NewPostUserLogin(username, password string) *PostUserLogin {
	return &PostUserLogin{Username: username, Password: password}
}

// 将用户数据持久化到数据库并检查是否出现用户名重复的现象
func (u *PostUserLogin) PersistData() error {
	return nil
}

func (u *PostUserLogin) SetToken() error {
	token, err := middleware.JwtGenerateToken(u, time.Hour)
	if err != nil {
		return err
	}
	u.Token = token
	return nil
}

// 检查用户登录时的用户名和密码是否合法
func (u *PostUserLogin) CheckPost() error {
	if u.Username == "" {
		return errors.New("用户名不能为空")
	}
	if len(u.Username) > MaxUsernameLength {
		return errors.New("用户名过长")
	}

	if len(u.Password) > MaxPasswordLength {
		return errors.New("密码过长")
	}
	if len(u.Password) < MinPasswordLength {
		return errors.New("密码不能少于5位")
	}
	return nil
}

// 返回正确信息
func RegisterOK(c *gin.Context, login *PostUserLogin) {
	c.JSON(http.StatusOK, UserRegisterReponse{
		CommonResponse: Model.CommonResponse{
			StatusCode: 0,
		},
		UserID: login.Userid,
		Token:  login.Token,
	})
}

// 返回错误信息
func RegisterErr(c *gin.Context, errmessage string) {
	c.JSON(http.StatusOK, UserRegisterReponse{
		CommonResponse: Model.CommonResponse{
			StatusCode: 1,
			StatusMsg:  errmessage,
		},
	})
}
