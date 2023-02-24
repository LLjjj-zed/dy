package middleware

import (
	"douyin.core/config"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type CommonResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// parseToken 解析token
func ParseToken(token string) (*jwt.StandardClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(config.Secret), nil
	})
	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
			return claim, nil
		}
	}
	return nil, err
}

func JWT() gin.HandlerFunc {
	return func(context *gin.Context) {
		//auth := context.Request.Header.Get("Authorization")
		auth := context.Query("token")
		if len(auth) == 0 {
			context.Abort()
			context.JSON(http.StatusUnauthorized, CommonResponse{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
		}
		auth = strings.Fields(auth)[1]
		token, err := ParseToken(auth)
		if err != nil {
			context.Abort()
			context.JSON(http.StatusUnauthorized, CommonResponse{
				StatusCode: -1,
				StatusMsg:  "Token Error",
			})
		} else {
			println("token 正确")
		}
		context.Set("userId", token.Id)
		context.Next()
	}
}

func JWTNOTOKEN() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.Query("token")
		var userId string
		if len(auth) == 0 {
			userId = "0"
		} else {
			auth = strings.Fields(auth)[1]
			token, err := ParseToken(auth)
			if err != nil {
				context.Abort()
				context.JSON(http.StatusUnauthorized, CommonResponse{
					StatusCode: -1,
					StatusMsg:  "Token Error",
				})
			} else {
				userId = token.Id
				println("token 正确")
			}
		}
		context.Set("userId", userId)
		context.Next()
	}
}

func JWTBody() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.Request.PostFormValue("token")
		fmt.Printf("%v \n", auth)

		if len(auth) == 0 {
			context.Abort()
			context.JSON(http.StatusUnauthorized, CommonResponse{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
		}
		auth = strings.Fields(auth)[1]
		token, err := ParseToken(auth)
		if err != nil {
			context.Abort()
			context.JSON(http.StatusUnauthorized, CommonResponse{
				StatusCode: -1,
				StatusMsg:  "Token Error",
			})
		} else {
			println("token 正确")
		}
		context.Set("userId", token.Id)
		context.Next()
	}
}
