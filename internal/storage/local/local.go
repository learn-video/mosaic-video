package local

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mauricioabreu/mosaic-video/internal/config"
)

type Client struct {
	cfg *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) CreateBucket(bucketName string) error {
	path := fmt.Sprintf("%s/%s", c.cfg.LocalStorage.Path, bucketName)

	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create mosaic directory path=%s  error=%w", path, err)
		}
	}

	return nil
}

func (c *Client) Upload(filename string, data []byte) error {
	return errors.New("not implemented")
}

func (c *Client) Get(filename string) (io.Reader, error) {
	path := fmt.Sprintf("%s/%s", c.cfg.LocalStorage.Path, filename)
	return os.Open(path)
}
