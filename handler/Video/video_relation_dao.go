package Video

// 用户视频关系表
type VideoRelation struct {
	VideoID int64 `json:"videoid"` //视频id
	UserID  int64 `json:"userid"`  //用户id
	Status  bool  `json:"status"`  //是否点赞
	Seen    bool  `json:"seen"`    //是否看过
}

// 视频列表
type VideoList struct {
	Videos []*Video
}

// 用户视频关系数据操作结构体
type VideoRelationDao struct {
}

// 用户视频关系数据操作结构体构造函数
func NewVideoRelationDao() *VideoRelationDao {
	return &VideoRelationDao{}
}
