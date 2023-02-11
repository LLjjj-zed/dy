package User

import "douyin.core/Model"

type UserRelation struct {
	UserID   int64 `json:"user_id"`   //用户id
	FollowID int64 `json:"follow_id"` //被关注者id
	Status   bool  `json:"status"`    //是否关注
}

type UserList struct {
	Users []*User
}

type UserRelationDao struct {
}

func NewUserRelationDao() *UserRelationDao {
	return &UserRelationDao{}
}

func (dao *UserRelationDao) QueryStatus(userid, followid int64) bool {
	var OK bool
	Model.DB.Select("status").Where("user_id=? AND folow_id=?", userid, followid).Find(&OK)
	return OK
}

func (dao *UserRelationDao) QueryFans(followid int64) *UserList {
	var users []*User
	Model.DB.Where("follow_id=?", followid).Find(&users)
	return &UserList{Users: users}
}

func (dao *UserRelationDao) QueryFollow(userid int64) *UserList {
	var users []*User
	Model.DB.Where("user_id=?", userid).Find(&users)
	return &UserList{Users: users}
}
