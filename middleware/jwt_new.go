package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Claims struct {
	UserId int64
	jwt.StandardClaims
}

var jwtKey = []byte("codegodrode")

// ReleaseToken 颁发token
func ReleaseToken(Userid int64) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: Userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "douyin-demo-seven-hulu",
			Subject:   "L_B__",
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析token
func ParseToken(tokenString string) (*Claims, bool) {
	token, _ := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if token != nil {
		if key, ok := token.Claims.(*Claims); ok {
			if token.Valid {
				return key, true
			} else {
				return key, false
			}
		}
	}
	return nil, false
}

// JWTMiddleWare 鉴权中间件，鉴权并设置user_id
func JWTMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetString("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}
		//用户不存在
		if tokenStr == "" {
			c.JSON(http.StatusOK, CommonResponse{StatusCode: 401, StatusMsg: "用户不存在"})
			c.Abort() //阻止执行
			return
		}
		//验证token
		UserClaims, ok := ParseToken(tokenStr)
		if !ok {
			c.JSON(http.StatusOK, CommonResponse{
				StatusCode: 403,
				StatusMsg:  "token不正确",
			})
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > UserClaims.ExpiresAt {
			c.JSON(http.StatusOK, CommonResponse{
				StatusCode: 402,
				StatusMsg:  "token过期",
			})
			c.Abort() //阻止执行
			return
		}
		c.Set("user_id", UserClaims.UserId)
		c.Next()
	}
}
