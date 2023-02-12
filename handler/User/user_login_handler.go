package User

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 用户登录回复结构体
type UserLoginResponse struct {
	CommonResponse
	Token  string `json:"token"`   // 用户鉴权token
	UserID int64  `json:"user_id"` // 用户id
}

// 用户登录回复结构体构造函数
func NewUserLoginResponse() *UserLoginResponse {
	return &UserLoginResponse{}
}

func UserLoginHandler(c *gin.Context) {
	//获取用户名和密码
	username := c.Query("username")
	get, exists := c.Get("password")
	if !exists {
		LoginErr(c, "密码不能为空")
		return
	}
	password := get.(string)
	//创建用户登录表对象
	userlogin := NewUserLoginTable(username, password)
	//创建用户登陆表数据操作对象
	userlogindao := NewUserRigisterDao()
	//验证用户账户密码
	err := userlogindao.QueryUserLogin(username, password, userlogin)
	if err != nil {
		LoginErr(c, err.Error())
		return
	}
	//获取token
	postUserLogin := NewPostUserLogin(username, password)
	err = postUserLogin.SetToken()
	if err != nil {
		LoginErr(c, err.Error())
		return
	}
	//赋值
	loginresponse := NewUserLoginResponse()
	loginresponse.UserID = userlogin.UserId
	loginresponse.Token = postUserLogin.Token
	LoginOK(c, loginresponse)
}

// 返回正确信息
func LoginOK(c *gin.Context, login *UserLoginResponse) {
	c.JSON(http.StatusOK, UserLoginResponse{
		CommonResponse: CommonResponse{
			StatusCode: 0,
		},
		UserID: login.UserID,
		Token:  login.Token,
	})
}

// 返回错误信息
func LoginErr(c *gin.Context, errmessage string) {
	c.JSON(http.StatusOK, UserLoginResponse{
		CommonResponse: CommonResponse{
			StatusCode: 1,
			StatusMsg:  errmessage,
		},
	})
}
