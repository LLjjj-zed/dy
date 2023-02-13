package Video

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserPublishListResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	VideoList  []*Video
}

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

func GetUserPublishList(userid int64) (*VideoList, error) {
	dao := NewVideoDao()
	list, err := dao.QueryUserPublishList(userid)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func UserPublishListOK(c *gin.Context, list *VideoList) {
	c.JSON(http.StatusOK, &UserPublishListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  list.Videos,
	})
}

func UserPublishListErr(c *gin.Context, errmessage string) {
	c.JSON(http.StatusOK, &UserPublishListResponse{
		StatusCode: 1,
		StatusMsg:  errmessage,
	})
}
