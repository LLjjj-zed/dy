package Interact

import "fmt"

const (
	FOLLOWED   int = 1
	UNFOLLOWED int = 2
)

type Relation struct {
	FollowCount   int64 `gorm:"follow_count" json:"follow_count"`     // 关注总数
	FollowerCount int64 `gorm:"follower_count" json:"follower_count"` // 粉丝总数
	ID            int64 `gorm:"id" json:"id"`                         // 用户id
}

func FollowCountkey(id int64) string {
	name := fmt.Sprintf("CommonUser%d'sFollowCount", id)
	return name
}
func FollowerCountkey(id int64) string {
	name := fmt.Sprintf("CommonUser%d'sFollowerCount", id)
	return name
}
func FollowSetkey(id string) string {
	name := fmt.Sprintf("CommonUser%d'sCount", id)
	return name
}
func FollowerSetkey(id string) string {
	name := fmt.Sprintf("CommonUser%d'sCount", id)
	return name
}
