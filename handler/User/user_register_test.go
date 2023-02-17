package User

import (
	"douyin.core/middleware"
	"fmt"
	"testing"
)

func TestSetToken(t *testing.T) {
	user := NewPostUserLogin("lljjj", "12345678")
	user.UserIdGenarate()
	fmt.Println(user.Userid)
	err := user.SetToken()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(user.Token)
	userclames, err := middleware.JwtParseUser(user.Token)
	if err != nil {
		t.Error(err)
	}
	if userclames.Userid != user.Userid {
		t.Error("error not equal")
	}
}

func TestCheckPost(t *testing.T) {
	users := []*PostUserLogin{}
	users = append(users, &PostUserLogin{Username: "", Password: "12345678"})
	users = append(users, &PostUserLogin{Username: "wfaesghretjsyktudlifukdyjthgrdsef", Password: "feawgsrhdtsjykdultyrh"})
	users = append(users, &PostUserLogin{Username: "fwaf", Password: "123"})
	users = append(users, &PostUserLogin{Username: "weafegarhreag", Password: "1fewasfaGrhbtjrshgerwger78"})
	users = append(users, &PostUserLogin{Username: "fawegwawh4", Password: "12gawer8"})
	users = append(users, &PostUserLogin{Username: "ewaf", Password: "1gaewdsfe8"})
	for _, user := range users {
		err := user.CheckPost()
		if err != nil {
			fmt.Println(err)
		}
	}
}
