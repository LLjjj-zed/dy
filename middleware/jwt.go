package middleware

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var AppSecret = ""
var AppIss = "https://github.com/132982317/douyin-demo"

type CommonResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

// 自定义payload结构体,不建议直接使用 dgrijalva/jwt-go `jwt.StandardClaims`结构体.因为他的payload包含的用户信息太少.
type UserStdClaims struct {
	jwt.StandardClaims
	Userid int64
}

// Valid 实现 `type Claims interface` 的 `Valid() error` 方法,自定义校验内容
func (c UserStdClaims) Valid() (err error) {
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

// JwtGenerateToken token生成
func JwtGenerateToken(Userid int64) (string, error) {
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expireTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        fmt.Sprintf("%d", Userid),
		Issuer:    AppIss,
	}

	uClaims := UserStdClaims{
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
func JwtParseUser(tokenString string) (*UserStdClaims, error) {
	if tokenString == "" {
		return nil, errors.New("no token is found in Authorization Bearer")
	}
	claims := UserStdClaims{}
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

func JWT() gin.HandlerFunc { //roleId uint
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		// 将 token 中的 claims 存入 context 中以便后续处理
		claims, err := JwtParseUser(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set("claims.Userid", claims.Userid)
		c.Next()
	}
}
