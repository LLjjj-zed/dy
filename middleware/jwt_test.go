package middleware

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestJwtGenerateToken(t *testing.T) {
	var userid int64 = 214312
	token, err := JwtGenerateToken(userid, time.Hour)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)
}

func TestJwtParseUser(t *testing.T) {
	var userid int64 = 214312
	token, err := JwtGenerateToken(userid, time.Hour)
	if err != nil {
		log.Fatal(err)
	}
	user, err := JwtParseUser(token)
	if err != nil {
		log.Fatal(err)
	}
	if user.Userid != userid {
		log.Fatal("id解析错误")
	}
}
