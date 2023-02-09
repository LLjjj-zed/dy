package user

import (
	"douyin.core/Model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserResponse struct {
	Model.CommonResponse
	User *User `json:"user"` // 用户信息
}

func UserInfoHandler(c *gin.Context) {
	id, exists := c.Get("user_id")
	if !exists {
		UserInfoErr(c, 3, "获取用户id失败")
		return
	}
	userDao := NewUserInfoDao()
	userinfo, err := userDao.GetUserByuserID(id)
	if err != nil {
		UserInfoErr(c, 4, "获取用户信息失败")
		return
	}
	UserInfoOK(c, userinfo)
}

// 返回正确信息
func UserInfoOK(c *gin.Context, login *User) {
	c.JSON(http.StatusOK, UserResponse{
		CommonResponse: Model.CommonResponse{
			StatusCode: 0,
		},
		User: &User{
			ID:            login.ID,
			Name:          login.Name,
			FollowerCount: login.FollowerCount,
			FollowCount:   login.FollowCount,
			IsFollow:      login.IsFollow,
		},
	})
}

// 返回错误信息
func UserInfoErr(c *gin.Context, code int64, errmessage string) {
	c.JSON(http.StatusOK, UserResponse{
		CommonResponse: Model.CommonResponse{
			StatusCode: code,
			StatusMsg:  errmessage,
		},
	})
}
