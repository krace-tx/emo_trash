package oss

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient 封装MinIO操作
type MinioClient struct {
	client     *minio.Client
	bucketName string
}

// Config MinIO配置
type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// NewMinioClient 创建MinIO客户端
func NewMinioClient(cfg Config) (*MinioClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// 检查桶是否存在，不存在则创建
	exists, err := client.BucketExists(context.Background(), cfg.Bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = client.MakeBucket(context.Background(), cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &MinioClient{
		client:     client,
		bucketName: cfg.Bucket,
	}, nil
}

// UploadFile 上传文件到MinIO
func (m *MinioClient) UploadFile(objectName string, reader io.Reader, objectSize int64) error {
	_, err := m.client.PutObject(context.Background(), m.bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	return err
}

// DownloadFile 从MinIO下载文件
func (m *MinioClient) DownloadFile(objectName string) (io.Reader, error) {
	return m.client.GetObject(context.Background(), m.bucketName, objectName, minio.GetObjectOptions{})
}

// ListFiles 列出文件
func (m *MinioClient) ListFiles(prefix string) ([]string, error) {
	var files []string
	objectCh := m.client.ListObjects(context.Background(), m.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		files = append(files, object.Key)
	}
	return files, nil
}

// GeneratePresignedURL 生成预签名URL
func (m *MinioClient) GeneratePresignedURL(objectName string, expiry time.Duration) (string, error) {
	u, err := m.client.PresignedGetObject(context.Background(), m.bucketName, objectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// DeleteFile 删除文件
func (m *MinioClient) DeleteFile(objectName string) error {
	return m.client.RemoveObject(context.Background(), m.bucketName, objectName, minio.RemoveObjectOptions{})
}
