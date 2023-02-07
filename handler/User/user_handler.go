package user

import "douyin.core/Model"

// User
type User struct {
	FollowCount   int64  `json:"follow_count"`  // 关注总数
	FollowerCount int64  `json:"follower_count"`// 粉丝总数
	ID            int64  `json:"id"`            // 用户id
	IsFollow      bool   `json:"is_follow"`     // true-已关注，false-未关注
	Name          string `json:"name"`          // 用户名称
}


type UserResponse struct {
	Model.CommonResponse
	User       *User   `json:"user"`       // 用户信息
}

func (u *User)NewUser()  {
	
}
