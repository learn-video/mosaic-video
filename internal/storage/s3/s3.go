package s3

import (
	"bytes"
	"context"
	"io"
	"path/filepath"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	client     *minio.Client
	bucketName string
}

func NewClient(cfg *config.Config) (storage.Storage, error) {
	client, err := minio.New(cfg.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3.AccessKeyID, cfg.S3.SecretAccessKey, ""),
		Secure: false,
	})

	return &Client{client: client, bucketName: cfg.S3.BucketName}, err
}

func (s3c *Client) CreateBucket(bucketName string) error {
	return s3c.client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
}

func (s3c *Client) Upload(filename string, data []byte) error {
	reader := bytes.NewReader(data)
	_, err := s3c.client.PutObject(
		context.Background(),
		s3c.bucketName,
		filename,
		reader,
		int64(len(data)),
		minio.PutObjectOptions{ContentType: getMIMEType(filename)},
	)

	return err
}

func (s3c *Client) Get(filename string) (io.Reader, error) {
	output, err := s3c.client.GetObject(
		context.Background(),
		s3c.bucketName,
		filename,
		minio.GetObjectOptions{},
	)

	if err != nil {
		return nil, err
	}

	return output, nil
}

func getMIMEType(path string) string {
	switch filepath.Ext(path) {
	case ".ts":
		return "video/mp2t"
	case ".m3u8":
		return "application/vnd.apple.mpegurl"
	default:
		return ""
	}
}
