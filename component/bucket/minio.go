package bucket

import (
	"bytes"
	"context"
	"io"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio interface {
	UploadFileFromBuffer(bucketName, objectName string, fileBuffer *bytes.Buffer, contentType string) error
	DownloadFileToBuffer(bucketName, objectName string) (*bytes.Buffer, error)
	DeleteFile(bucketName, objectName string) error
}

type MinioClient struct {
	client *minio.Client
}

func NewMinioClient() *MinioClient {
	endpoint := "localhost:9000"
	accessKeyID := "admin"
	secretAccessKey := "admin123"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &MinioClient{client: minioClient}
}

func (mc *MinioClient) UploadFileFromBuffer(bucketName, objectName string, fileBuffer *bytes.Buffer, contentType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// 检查存储桶是否存在
	exists, err := mc.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		// 如果桶不存在则创建桶
		err := mc.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	reader := bytes.NewReader(fileBuffer.Bytes())

	_, err = mc.client.PutObject(ctx, bucketName, objectName, reader, int64(fileBuffer.Len()), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return err
	}

	return nil
}

func (mc *MinioClient) DownloadFileToBuffer(bucketName, objectName string) (*bytes.Buffer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// 从 MinIO 获取对象
	object, err := mc.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, object)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (mc *MinioClient) DeleteFile(bucketName, objectName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := mc.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
