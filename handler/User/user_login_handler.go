package User

import (
	"douyin.core/Model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserLoginResponse 用户登录回复结构体
type UserLoginResponse struct {
	CommonResponse
	Token  string `json:"token"`   // 用户鉴权token
	UserID int64  `json:"user_id"` // 用户id
}

// NewUserLoginResponse 用户登录回复结构体构造函数
func NewUserLoginResponse() *UserLoginResponse {
	return &UserLoginResponse{}
}

// UserLoginHandler 用户登录处理函数
func UserLoginHandler(c *gin.Context) {
	//获取用户名和密码
	username := c.Query("username")
	password, exists := c.GetQuery("password")
	if !exists {
		LoginErr(c, "密码不能为空")
		return
	}
	//创建用户登录表对象
	userlogin := Model.NewUserLoginTable(username, password)
	//创建用户登陆表数据操作对象
	userlogindao := Model.NewUserRigisterDao()
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
	LoginOK(c, userlogin.UserId, postUserLogin.Token)
}

// LoginOK 返回正确信息
func LoginOK(c *gin.Context, userid int64, token string) {
	c.JSON(http.StatusOK, UserLoginResponse{
		CommonResponse: CommonResponse{
			StatusCode: 0,
		},
		UserID: userid,
		Token:  token,
	})
}

// LoginErr 返回错误信息
func LoginErr(c *gin.Context, errmessage string) {
	c.JSON(http.StatusOK, UserLoginResponse{
		CommonResponse: CommonResponse{
			StatusCode: 1,
			StatusMsg:  errmessage,
		},
	})
}
