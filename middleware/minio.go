package middleware

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func Initminio() *minio.Client {
	endpoint := "8.130.73.119:9000"
	accessKeyID := "douyin"
	secretAccessKey := "88888888"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now set up
	return minioClient
}

func UploadVideoToMinio(minioClient *minio.Client, videoname, videopath, bucketName string) error {
	// Upload the mp4 file with FPutObject
	info, err := minioClient.FPutObject(context.Background(), bucketName, videoname, videopath, minio.PutObjectOptions{ContentType: "video/mp4"})
	if err != nil {
		return err
	}
	log.Printf("Successfully uploaded %s of size %d\n", videoname, info.Size)
	return nil
}

func UploadImageoMinio(minioClient *minio.Client, imagename, imagepath, bucketName string) error {
	// Upload the png file with FPutObject
	info, err := minioClient.FPutObject(context.Background(), bucketName, imagename, imagepath, minio.PutObjectOptions{ContentType: "image/jpg"})
	if err != nil {
		return err
	}
	log.Printf("Successfully uploaded %s of size %d\n", imagename, info.Size)
	return nil
}
