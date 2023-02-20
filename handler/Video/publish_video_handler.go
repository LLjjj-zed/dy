package Video

import (
	"douyin.core/Model"
	"douyin.core/middleware"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	minioClient    *minio.Client
	NewminioClient sync.Once
)

// 下载文件的协程数量
const numWorkers = 4

// GetminioClient 连接minio客户端
func GetminioClient() {
	NewminioClient.Do(func() {
		minioClient = middleware.Initminio()
	})
}

// PublishVideoResponse 用户投稿回复结构体
type PublishVideoResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

// PublishVedioHandler 用户投稿处理函数，用于处理http请求
func PublishVedioHandler(c *gin.Context) {
	//从请求中获取视频标题和token
	title := c.PostForm("title")
	token, ok := c.GetQuery("token")
	if !ok {
		PublishVideoErr(c, "未能成功获取token，请重试")
		return
	}
	//从token中解析出用户id
	userclaim, err := middleware.JwtParseUser(token)
	if err != nil {
		PublishVideoErr(c, err.Error())
	}
	userid := userclaim.Userid
	//从请求中获取时视频数据
	file, err := c.FormFile("data")
	if err != nil {
		PublishVideoErr(c, err.Error())
		return
	}
	//从视频数据中获取视频后缀格式
	ext := filepath.Ext(file.Filename)
	if ext != ".mp4" {
		PublishVideoErr(c, "上传视频格式错误，请重试")
		return
	}
	videoDao := Model.NewVideoDao()
	//获取用户视频序号，用于生成视频文件名
	codeint, err := videoDao.GetUserVideoCode(userid)
	if err != nil {
		PublishVideoErr(c, err.Error())
		return
	}
	//将用户视频序号转换成字符串
	code := strconv.Itoa(int(codeint))
	var userinfo Model.UserInfoDao
	//获取用户名，用于生成视频文件名
	username, err := userinfo.GetUserNameByUserID(userid)
	name := username
	if err != nil {
		PublishVideoErr(c, err.Error())
	}
	//将视频持久化到本地，使用strings.Builder替换+提高性能
	videoname := GetFilename(name, code, ext)
	path := filepath.Join("./public/", videoname)
	dst, err := os.Create(path)
	defer dst.Close()
	if err != nil {
		PublishVideoErr(c, err.Error())
		return
	}
	//多协程下载文件 todo
	stat := DownloadFile(c, file, dst)
	// 当多协程下载失败时进行降级
	if stat != "nil" {
		err = c.SaveUploadedFile(file, path)
		if err != nil {
			PublishVideoErr(c, err.Error())
			return
		}
	}
	imagename := GetFilename(name, code, ".jpg")
	//生成截图
	err = middleware.GetSnapshotCmd(videoname, imagename)
	if err != nil {
		PublishVideoErr(c, err.Error())
		return
	}
	//连接minio
	GetminioClient()
	//将视频上传至minio
	err = middleware.UploadVideoToMinio(c, minioClient, videoname, path, "video")
	if err != nil {
		PublishVideoErr(c, err.Error())
	}
	//将视频信息持久化到数据库
	err = videoDao.PersistNewVideo(title, userid, codeint, videoname, imagename, &userinfo)
	if err != nil {
		PublishVideoErr(c, err.Error())
	}
	PublishVideoOk(c)
}

// GetFilename 获取文件名，根据用户名，用户视频序号
func GetFilename(name, code, ext string) string {
	var build strings.Builder
	build.WriteString(name)
	build.WriteString("-code-")
	build.WriteString(code)
	build.WriteString(ext)
	filename := build.String()
	return filename
}

func DownloadFile(c *gin.Context, file *multipart.FileHeader, dst *os.File) (Err string) {
	// 获取上传文件的大小
	size := file.Size

	// 计算每个协程需要下载的字节数
	chunkSize := size / numWorkers

	// 等待协程完成的计数器
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	var download func(int64, int64)

	// 启动协程下载字节范围内的数据
	download = func(start, end int64) {
		defer wg.Done()

		// 创建HTTP请求
		req, err := http.NewRequest("GET", c.Request.RequestURI, nil)
		if err != nil {
			Err = err.Error()
		}

		// 添加Range头
		req.Header.Add("Range", "bytes="+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10))

		// 发送请求
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			Err = err.Error()
		}
		defer resp.Body.Close()

		// 复制数据到目标文件
		_, err = io.Copy(dst, resp.Body)
		if err != nil {
			Err = err.Error()
		}
	}

	// 启动多个协程下载文件
	for i := 0; i < numWorkers; i++ {
		// 计算当前协程需要下载的字节范围
		start := int64(i) * chunkSize
		end := start + chunkSize - 1
		if i == numWorkers-1 {
			end = size - 1
		}
		go download(start, end)
	}

	// 等待所有协程完成
	wg.Wait()
	return "nil"
}

// PublishVideoOk 返回正确信息
func PublishVideoOk(c *gin.Context) {
	c.JSON(http.StatusOK, &PublishVideoResponse{
		StatusCode: 0,
	})
}

// PublishVideoErr  返回错误信息
func PublishVideoErr(c *gin.Context, errmeassage string) {
	c.JSON(http.StatusOK, &PublishVideoResponse{
		StatusCode: 1,
		StatusMsg:  errmeassage,
	})
}
