package Video

import "douyin.core/Model"

type VideoRelation struct {
	VideoID int64 `json:"videoid"` //视频id
	UserID  int64 `json:"userid"`  //用户id
	Status  bool  `json:"status"`  //是否点赞
}

type VideoList struct {
	Videos []*Video
}

type VideoRelationDao struct {
}

func NewVideoRelationDao() *VideoRelationDao {
	return &VideoRelationDao{}
}

func (dao *VideoRelationDao) QueryStatus(userid, videoid int64) bool {
	var OK bool
	Model.DB.Select("status").Where("userid=? AND videoid=?", userid, videoid).Find(&OK)
	return OK
}
