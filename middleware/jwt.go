package middleware

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

var AppSecret = "" //viper.GetString会设置这个值(32byte长度)
var AppIss = "https://github.com/132982317/douyin-demo"

type CommonResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

// 自定义payload结构体,不建议直接使用 dgrijalva/jwt-go `jwt.StandardClaims`结构体.因为他的payload包含的用户信息太少.
type userStdClaims struct {
	jwt.StandardClaims
	Userid int64
}

// 实现 `type Claims interface` 的 `Valid() error` 方法,自定义校验内容
func (c userStdClaims) Valid() (err error) {
	//token过期校验
	if c.VerifyExpiresAt(time.Now().Unix(), true) == false {
		return errors.New("token is expired")
	}
	//令牌发行者校验
	if !c.VerifyIssuer(AppIss, true) {
		return errors.New("token's issuer is wrong")
	}
	if c.Userid < 1 {
		return errors.New("invalid user in jwt")
	}
	return
}

// token生成
func JwtGenerateToken(Userid int64, d time.Duration) (string, error) {
	expireTime := time.Now().Add(d)
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expireTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        fmt.Sprintf("%d", Userid),
		Issuer:    AppIss,
	}

	uClaims := userStdClaims{
		StandardClaims: stdClaims,
		Userid:         Userid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uClaims)
	// 使用密钥对完整的编码令牌进行签名并获取字符串形式
	tokenString, err := token.SignedString([]byte(AppSecret))
	if err != nil {
		logrus.WithError(err).Fatal("config is wrong, can not generate jwt")
	}
	return tokenString, err
}

// JwtParseUser 解析payload的内容,得到用户信息
// gin-middleware 会使用这个方法
func JwtParseUser(tokenString string) (*userStdClaims, error) {
	if tokenString == "" {
		return nil, errors.New("no token is found in Authorization Bearer")
	}
	claims := userStdClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(AppSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return &claims, err
}

const contextKeyUserObj = "authedUserObj"
const bearerLength = len("Bearer ")

func ctxTokenToUser(c *gin.Context) { //roleId uint
	token, ok := c.GetQuery("token")
	if !ok {
		hToken := c.GetHeader("Authorization")
		if len(hToken) < bearerLength {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": "header Authorization has not Bearer token"})
			return
		}
		token = strings.TrimSpace(hToken[bearerLength:])
	}
	usr, err := JwtParseUser(token)
	if err != nil {
		c.JSON(http.StatusOK, CommonResponse{
			StatusCode: 403,
			StatusMsg:  "token不正确",
		})
		c.Abort() //阻止执行
		return
	}

	//账号权限
	//if (usr.RoleId & roleId) != roleId {
	//	c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": "roleId 没有权限"})
	//	return
	//}

	//token超时
	if time.Now().Unix() > usr.ExpiresAt {
		c.JSON(http.StatusOK, CommonResponse{
			StatusCode: 402,
			StatusMsg:  "token过期",
		})
		c.Abort() //阻止执行
		return
	}

	//store the user Model in the context
	c.Set(contextKeyUserObj, *usr)
	c.Next()
	// after request
}
