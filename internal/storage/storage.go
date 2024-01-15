package storage

import "io"

type Storage interface {
	CreateBucket(bucketName string) error
	Upload(path string, data []byte) error
	Get(path string) (io.Reader, error)
}
