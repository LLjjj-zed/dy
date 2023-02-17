package service

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

type MinioConfig struct {
	endpoint        string //存储服务URL
	accessKeyID     string //用户ID
	secretAccessKey string //账户密码
	secure          bool   //是否使用https
}

var cfg = MinioConfig{
	endpoint:        "localhost:9000",
	accessKeyID:     "jerry",
	secretAccessKey: "88888888",
	secure:          false,
}

var (
	client *minio.Client
	err    error
)

func init() {
	client, err := minio.New(cfg.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.accessKeyID, cfg.secretAccessKey, ""),
		Secure: cfg.secure,
	})
	if err != nil {
		log.Fatalf("Error connecting to MinIO service: %s", err)
	}
	log.Printf("%#v\n", client)
}

func UploadObj(ctx context.Context) {
	uploadInfo, err := client.FPutObject(context.Background(), "douyin", "videoid", "videoid.mp4", minio.PutObjectOptions{
		ContentType: "video/mpeg4",
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Successfully uploaded:", uploadInfo)
}
