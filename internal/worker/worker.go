package worker

import (
	"context"
	"log"
	"time"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/locking"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic/command"
	"github.com/mauricioabreu/mosaic-video/internal/storage"
)

const (
	LockingTimeTTL    time.Duration = 120 * time.Second
	KeepAliveInterval time.Duration = 30 * time.Second
)

func GenerateMosaic(m mosaic.Mosaic, cfg *config.Config, locker locking.Locker, cmdExecutor mosaic.Command, runningProcesses map[string]bool, stg storage.Storage) error {
	_, exists := runningProcesses[m.Name]
	if exists {
		return nil
	}

	if err := createBucket(&m, cfg, stg); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lock, err := locker.Obtain(ctx, m.Name, LockingTimeTTL)
	if err != nil {
		return err
	}

	go keepAlive(ctx, lock)

	args := command.Build(m, cfg)
	if err := cmdExecutor.Execute("ffmpeg", args...); err != nil {
		if lerr := lock.Release(ctx); lerr != nil {
			return lerr
		}

		return err
	}

	runningProcesses[m.Name] = true

	return nil
}

func createBucket(m *mosaic.Mosaic, cfg *config.Config, stg storage.Storage) error {
	if cfg.StorageType.IsLocal() {
		return stg.CreateBucket(m.Name)
	}

	return stg.CreateBucket(cfg.S3.BucketName)
}

func keepAlive(ctx context.Context, lock locking.Lock) {
	ticker := time.NewTicker(KeepAliveInterval)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			if err := lock.Refresh(ctx, LockingTimeTTL); err != nil {
				log.Println("failed to refresh lock TTL, error=%w", err)
			}
		}
	}
}
