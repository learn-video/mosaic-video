package storage

type Storage interface {
	CreateBucket(bucketName string) error
	Upload(path string, data []byte) error
}
