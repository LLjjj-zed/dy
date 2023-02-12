package User

import "douyin.core/Model"

// 用户关系表
type UserRelation struct {
	UserID   int64 `json:"user_id"`   //用户id
	FollowID int64 `json:"follow_id"` //被关注者id
	Status   bool  `json:"status"`    //是否关注
}

// 用户列表，用于返回用户列表
type UserList struct {
	Users []*User
}

// 用户关系表数据操作结构体
type UserRelationDao struct {
}

// 用户关系表数据操作结构体构造函数
func NewUserRelationDao() *UserRelationDao {
	return &UserRelationDao{}
}

// 查询用户关系，场景：刷视频的时候获取用户信息查看视频发布者是否已经关注
func (dao *UserRelationDao) QueryStatus(userid, followid int64) bool {
	var OK bool
	Model.DB.Select("status").Where("user_id=? AND folow_id=?", userid, followid).Find(&OK)
	return OK
}

// 查询粉丝
func (dao *UserRelationDao) QueryFans(followid int64) *UserList {
	var users []*User
	Model.DB.Where("follow_id=?", followid).Find(&users)
	return &UserList{Users: users}
}

// 查询关注
func (dao *UserRelationDao) QueryFollow(userid int64) *UserList {
	var users []*User
	Model.DB.Where("user_id=?", userid).Find(&users)
	return &UserList{Users: users}
}
