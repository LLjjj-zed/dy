package Video

import (
	"douyin.core/Model"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// UserPublishListResponse 发布列表回复结构体
type UserPublishListResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	VideoList  []*Model.Video
}

// UserPublishListHandler 发布列表处理函数
func UserPublishListHandler(c *gin.Context) {
	useridstr, exist := c.GetQuery("user_id")
	if exist {
		userid, err := strconv.ParseInt(useridstr, 10, 64)
		if err != nil {
			UserPublishListErr(c, err.Error())
			return
		}
		videoList, err := GetUserPublishList(userid)
		if err != nil {
			UserPublishListErr(c, err.Error())
			return
		}
		UserPublishListOK(c, videoList)
	}
	UserPublishListErr(c, errors.New("userid获取失败，请重试").Error())
}

// GetUserPublishList 获取发布列表
func GetUserPublishList(userid int64) (*Model.VideoList, error) {
	dao := Model.NewVideoDao()
	list, err := dao.QueryUserPublishList(userid)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// UserPublishListOK 返回正确信息
func UserPublishListOK(c *gin.Context, list *Model.VideoList) {
	c.JSON(http.StatusOK, &UserPublishListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  list.Videos,
	})
}

// UserPublishListErr 返回错误信息
func UserPublishListErr(c *gin.Context, errmessage string) {
	c.JSON(http.StatusOK, &UserPublishListResponse{
		StatusCode: 1,
		StatusMsg:  errmessage,
	})
}
