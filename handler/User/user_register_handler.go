package User

import (
	"douyin.core/Model"
	"douyin.core/config"
	"douyin.core/middleware"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

const (
	MaxUsernameLength = 20
	MaxPasswordLength = 18
	MinPasswordLength = 5
)

// RegisterReponse 用户注册回复结构体
type RegisterReponse struct {
	CommonResponse
	Token  string `json:"token"`   // 用户鉴权token
	UserID int64  `json:"user_id"` // 用户id
}

// PostUserLogin 用户登录校验
type PostUserLogin struct {
	Username string
	Password string
	Userid   int64
	Token    string
}

// UserRegistHandler 用户注册处理函数
func UserRegistHandler(c *gin.Context) {
	//获取用户名和密码
	username := c.Query("username")
	password, exist := c.GetQuery("password")
	if !exist {
		RegisterErr(c, errors.New("密码不能为空").Error())
	}
	//创建新对象
	newPostUserLogin := NewPostUserLogin(username, password)
	//注册新用户
	err := newPostUserLogin.Register()
	if err != nil {
		//返回错误信息
		RegisterErr(c, err.Error())
		return
	}
	//返回正确信息
	RegisterOK(c, newPostUserLogin)
}

// Register 注册新用户
func (u *PostUserLogin) Register() error {
	//校验参数
	err := u.CheckPost()
	if err != nil {
		return err
	}

	//生成userid
	u.UserIdGenarate()

	//持久化到数据库
	err = u.PersistData()
	if err != nil {
		return err
	}

	//获取token
	token, err := u.NewToken()
	if err != nil {
		return err
	}
	u.Token = token
	return nil
}

// NewPostUserLogin 创建对象，用于注册
func NewPostUserLogin(username, password string) *PostUserLogin {
	return &PostUserLogin{Username: username, Password: password}
}

// PersistData 将用户数据持久化到数据库并检查是否出现用户名重复的现象
func (u *PostUserLogin) PersistData() error {
	//创建用户表数据操作对象
	userDAO := Model.NewUserInfoDao()
	//创建用户注册表数据操作对象
	userRigestDao := Model.NewUserRigisterDao()
	//检查用户名是否已经存在
	_, err := userDAO.GetUserByUserName(u.Username)
	is := errors.Is(err, gorm.ErrRecordNotFound)
	if is {
		//将数据持久化到用户表
		err = userDAO.InsertToUserInfoTable(u.Userid, u.Username)
		if err != nil {
			return err
		}
		//对用户密码进行AES-256加密，保障用户安全
		passwordstr := []byte(u.Password)
		aes, err := middleware.EnPwdCode(passwordstr)
		if err != nil {
			return err
		}
		u.Password = aes
		//将数据持久化到用户注册表
		err = userRigestDao.RegistUsertoDb(u.Userid, u.Username, u.Password)
		if err != nil {
			return err
		}
		return nil
	} else {
		err := errors.New("用户名已存在，请重试")
		return err
	}
}

//// SetToken 获取token
//func (u *PostUserLogin) SetToken() error {
//	var login PostUserLogin
//	token, err := middleware.ReleaseToken(middleware.Login(login))
//	if err != nil {
//		return err
//	}
//	u.Token = token
//	return nil
//}

// UserIdGenarate 用户id生成
func (u *PostUserLogin) UserIdGenarate() {
	worker, _ := middleware.NewWorker(1)
	id := worker.GetId()
	u.Userid = id
}

// CheckPost 检查用户登录时的用户名和密码是否合法
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

// NewToken 根据信息创建token
func (u *PostUserLogin) NewToken() (string, error) {
	expiresTime := time.Now().Unix() + int64(7*24*time.Hour)
	fmt.Printf("expiresTime: %v\n", expiresTime)
	id64 := u.Userid
	fmt.Printf("id: %v\n", strconv.FormatInt(id64, 10))
	claims := jwt.StandardClaims{
		Audience:  u.Username,
		ExpiresAt: expiresTime,
		Id:        strconv.FormatInt(id64, 10),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "tiktok",
		NotBefore: time.Now().Unix(),
		Subject:   "token",
	}
	var jwtSecret = []byte(config.Secret)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token, err := tokenClaims.SignedString(jwtSecret); err == nil {
		token = "Bearer " + token
		println("generate token success!\n")
		return token, nil
	} else {
		println("generate token fail\n")
		return "fail", err
	}
}

// RegisterOK 返回正确信息
func RegisterOK(c *gin.Context, login *PostUserLogin) {
	c.JSON(http.StatusOK, RegisterReponse{
		CommonResponse: CommonResponse{
			StatusCode: 0,
		},
		UserID: login.Userid,
		Token:  login.Token,
	})
}

// RegisterErr 返回错误信息
func RegisterErr(c *gin.Context, errmessage string) {
	c.JSON(http.StatusOK, RegisterReponse{
		CommonResponse: CommonResponse{
			StatusCode: 1,
			StatusMsg:  errmessage,
		},
	})
}
