package middleware

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

var bucketName = "douyin"

func Initminio() (*minio.Client, context.Context) {
	endpoint := "http://23.94.57.209:9001"
	accessKeyID := "IlOzGpW6BAYRuKGO"
	secretAccessKey := "IzLLsT2WsN8xY8SymTtLu6pzMJDwNxpI"
	useSSL := true

	// Initialize minio client object.
	ctx := context.Background()
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now set up
	policy := "{\"Version\": \"2012-10-17\",\"Statement\": [{\"Effect\": \"Allow\",\"Action\": [\"s3:GetBucketLocation\"],\"Resource\": [\"arn:aws:s3:::starcity\"]},{\"Effect\": \"Allow\",\"Action\": [\"s3:ListBucket\"],\"Resource\": [\"arn:aws:s3:::starcity\"],\"Condition\": {\"StringEquals\": {\"s3:prefix\": [\"bear\",\"prefix/\"]}}},{\"Effect\": \"Allow\",\"Action\": [\"s3:GetObject\"],\"Resource\": [\"arn:aws:s3:::starcity/bear*\",\"arn:aws:s3:::starcity/prefix/*\"]}]}"
	minioClient.SetBucketPolicy(ctx, "douyin", policy)

	return minioClient, ctx
}

func UploadFileToMinio(minioClient *minio.Client, videoname, videopath string, ctx context.Context) {
	// Upload the mp4 file

	objectName := videoname
	filePath := videopath
	contentType := "video/mp4"

	// Upload the mp4 file with FPutObject
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}
