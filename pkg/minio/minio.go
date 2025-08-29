package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

var (
	MinioClient *minio.Client
	bucketName  = "douyin-videos"
)

func InitMinio() {
	endpoint := "localhost:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	// 初始化MinIO客户端
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}

	MinioClient = client

	// 确保存储桶存在
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("Failed to check bucket existence: %v", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
		}
	}
}

// GetVideoURL 获取视频URL
func GetVideoURL(objectName string) (string, error) {
	// 生成预签名URL，有效期1小时
	presignedURL, err := MinioClient.PresignedGetObject(context.Background(), bucketName, objectName, 3600, nil)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// UploadVideo 上传视频
func UploadVideo(objectName string, filePath string) error {
	_, err := MinioClient.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{})
	return err
}

// GetVideoChunk 获取视频分片
func GetVideoChunk(objectName string, offset int64, length int64) ([]byte, error) {
	opts := minio.GetObjectOptions{}
	opts.SetRange(offset, offset+length-1)

	object, err := MinioClient.GetObject(context.Background(), bucketName, objectName, opts)
	if err != nil {
		return nil, err
	}
	defer object.Close()

	buffer := make([]byte, length)
	_, err = object.Read(buffer)
	return buffer, err
}
