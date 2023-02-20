package Model

import (
	"douyin.core/dal"
	"douyin.core/handler/Interact"
	"strconv"
)

var conn = dal.Redisclient
var ctx = dal.Ctx

//前者关注后者，前者关注数量+1，后者粉丝数量+1,
//还要更新关注列表，粉丝列表
func Follow(userid, id int64) error {
	var key = string(rune(userid))
	var val = string(rune(id))
	pipe := conn.Pipeline()
	if exi, err := conn.SIsMember(ctx, Interact.FollowSetkey(key), val).Result(); exi != true {
		if err != nil {
			return err
		}
		//关注操作，前者的关注集合添加后者，后者的粉丝集合添加前者
		_, adderr := pipe.SAdd(ctx, Interact.FollowSetkey(key), val).Result()
		_, adderr = pipe.SAdd(ctx, Interact.FollowerSetkey(val), key).Result()
		if adderr != nil {
			return adderr
		}
	} else {
		return nil
	}
	//前者的关注数以及后者的粉丝数加1
	if follow := pipe.Incr(ctx, Interact.FollowCountkey(userid)); follow.Err() != nil {
		return follow.Err()
	} else {
		if followed := pipe.Incr(ctx, Interact.FollowerCountkey(id)); followed.Err() != nil {
			return follow.Err()
		}
	}
	_, err := pipe.Exec(ctx)
	if err != nil { // 报错后进行一次额外尝试
		_, err = pipe.Exec(ctx)
		if err != nil {
			return nil
		}
	}
	return nil
}

//前者取关后者，前者关注数量-1，后者粉丝数量-1
func UnFollow(userid, id int64) error {
	var key = string(rune(userid))
	var val = string(rune(id))
	pipe := conn.Pipeline()
	if exi, err := conn.SIsMember(ctx, Interact.FollowSetkey(key), val).Result(); exi == true {
		if err != nil {
			return err
		}
		_, adderr := pipe.SRem(ctx, Interact.FollowSetkey(key), val).Result() //val==1，关注操作，前者的关注集合添加后者，后者的粉丝集合添加前者
		_, adderr = pipe.SRem(ctx, Interact.FollowerSetkey(val), key).Result()
		if adderr != nil {
			return adderr
		}
	} else {
		return nil
	}
	//前者的关注数以及后者的粉丝数减1
	if follow := pipe.Decr(ctx, Interact.FollowCountkey(userid)); follow.Err() != nil {
		return follow.Err()
	} else {
		if followed := pipe.Decr(ctx, Interact.FollowerCountkey(id)); followed.Err() != nil {
			return follow.Err()
		}
	}
	_, err := pipe.Exec(ctx)
	if err != nil { // 报错后进行一次额外尝试
		_, err = pipe.Exec(ctx)
		if err != nil {
			return nil
		}
	}
	return nil
}
func GetFollowList(id int64) ([]User, error) {
	var user []User
	vals, getallerr := conn.SMembers(ctx, Interact.FollowSetkey(strconv.FormatInt(id, 10))).Result()
	if getallerr != nil {
		return user, getallerr
	}
	err := DB.Where("id in (?)", vals).Find(&user).Error
	return user, err
}
func GetFollowerList(id int64) ([]User, error) {
	var user []User
	vals, getallerr := conn.SMembers(ctx, Interact.FollowerSetkey(strconv.FormatInt(id, 10))).Result()
	if getallerr != nil {
		return user, getallerr
	}
	err := DB.Where("id in (?)", vals).Find(&user).Error
	return user, err
}

//查询双向关注的好友
func GetFriendsList(id int64) ([]User, error) {
	var user []User
	var friends []string
	vals, getallerr := conn.SMembers(ctx, Interact.FollowerSetkey(strconv.FormatInt(id, 10))).Result()
	if getallerr != nil {
		return user, getallerr
	}
	for _, val := range vals {
		if exi, err := conn.SIsMember(ctx, Interact.FollowerSetkey(val), id).Result(); err != nil {
			return user, err
		} else {
			if exi {
				friends = append(friends, val)
			}
		}
	}

	err := DB.Where("id in (?)", friends).Find(&user).Error
	return user, err
}
