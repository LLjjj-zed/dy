package social

import (
	"douyin.core/Model"
	"douyin.core/dal"
	"douyin.core/handler/Interact"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []Model.User `json:"user_list"`
}

var ctx = dal.Ctx
var DB = dal.DB
var conn = dal.Redisclient

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	action_type, _ := strconv.Atoi(c.Query("action_type"))
	id, _ := strconv.Atoi(c.Query("to_user_id"))
	user, exist := usersLoginInfo[token]
	if exist {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户不存在"})
	}
	pipe := conn.Pipeline()
	var key = string(rune(user.Id))
	var val = string(rune(id))
	var err error

	switch action_type {
	case Interact.FOLLOWED:
		{
			if exi, _ := pipe.SIsMember(ctx, Interact.FollowSetkey(key), val).Result(); exi != true {
				//关注操作，前者的关注集合添加后者，后者的粉丝集合添加前者
				if _, adderr := pipe.SAdd(ctx, Interact.FollowSetkey(key), val).Result(); adderr != nil {
					logrus.Info(adderr)
				}
				if _, adderr := pipe.SAdd(ctx, Interact.FollowerSetkey(val), key).Result(); adderr != nil {
					logrus.Info(adderr)
				}
				//前者的关注数以及后者的粉丝数加1
				if follow := pipe.Incr(ctx, Interact.FollowCountkey(user.Id)); follow.Err() != nil {
					logrus.Info(follow.Err())
				} else {
					if followed := pipe.Incr(ctx, Interact.FollowerCountkey(int64(id))); followed.Err() != nil {
						logrus.Info(follow.Err())
					}
				}
				_, err = pipe.Exec(ctx)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "关注操作失败"})
				}
			}
			if err == nil {
				c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "成功关注用户"})
			}
		}
	case Interact.UNFOLLOWED:
		{
			if exi, _ := pipe.SIsMember(ctx, Interact.FollowSetkey(key), val).Result(); exi == true {
				//关注操作，前者的关注集合添加后者，后者的粉丝集合添加前者
				if _, adderr := pipe.SRem(ctx, Interact.FollowSetkey(key), val).Result(); adderr != nil {
					logrus.Info(adderr)
				}
			}
			if exi, _ := pipe.SIsMember(ctx, Interact.FollowerSetkey(val), key).Result(); exi == true {
				//关注操作，前者的关注集合添加后者，后者的粉丝集合添加前者
				if _, adderr := pipe.SRem(ctx, Interact.FollowerSetkey(val), key).Result(); adderr != nil {
					logrus.Info(adderr)
				}
			}
			//前者的关注数以及后者的粉丝数加1
			if follow := pipe.Decr(ctx, Interact.FollowCountkey(user.Id)); follow.Err() != nil {
				logrus.Info(follow.Err())
			}
			if followed := pipe.Decr(ctx, Interact.FollowerCountkey(int64(id))); followed.Err() != nil {
				logrus.Info(followed.Err())
			}
			_, err = pipe.Exec(ctx)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "取消关注用户失败"})
			}
		}
		if err == nil {
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "取消关注用户成功"})
		}
		break
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	token := c.Query("token")
	id, _ := strconv.Atoi(c.Query("user_id"))
	_, exi := usersLoginInfo[token]
	if exi != true {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户不存在"})
	}
	var user []Model.User
	vals, getallerr := conn.SMembers(ctx, Interact.FollowSetkey(strconv.FormatInt(int64(id), 10))).Result()
	if getallerr != nil {
		c.AbortWithStatusJSON(200, getallerr)
	}
	err := DB.Where("id in (?)", vals).Find(&user).Error
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取关注列表失败"})
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserList: user,
		})
	}
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	_, exi := usersLoginInfo[token]
	id, _ := strconv.Atoi(c.Query("user_id"))
	if exi != true {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	var user []Model.User
	vals, getallerr := conn.SMembers(ctx, Interact.FollowerSetkey(strconv.FormatInt(int64(id), 10))).Result()
	if getallerr != nil {
		c.AbortWithStatusJSON(http.StatusOK, Response{1, "获取粉丝列表失败"})
	}
	err := DB.Where("id in (?)", vals).Find(&user).Error
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取粉丝列表失败"})
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserList: user,
		})
	}
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	token := c.Query("token")
	_, exi := usersLoginInfo[token]
	id, _ := strconv.Atoi(c.Query("user_id"))
	if exi != true {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户不存在"})
	}
	var err error
	var user []Model.User
	var friends []string
	vals, getallerr := conn.SMembers(ctx, Interact.FollowerSetkey(strconv.FormatInt(int64(id), 10))).Result()
	if getallerr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取好友列表失败"})
	}
	for _, val := range vals {
		if exi, err := conn.SIsMember(ctx, Interact.FollowerSetkey(val), id).Result(); err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取好友列表失败"})
		} else {
			if exi {
				friends = append(friends, val)
			}
		}
	}
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取好友列表失败"})
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserList: user,
		})
	}
}
