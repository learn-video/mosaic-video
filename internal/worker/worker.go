package worker

import (
	"context"
	"time"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/locking"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic/command"
	"github.com/mauricioabreu/mosaic-video/internal/storage"
)

const LockingTimeTTL time.Duration = 120 * time.Second

func GenerateMosaic(m mosaic.Mosaic, cfg *config.Config, locker locking.Locker, cmdExecutor mosaic.Command, runningProcesses map[string]bool, stg storage.Storage) error {
	_, exists := runningProcesses[m.Name]
	if exists {
		return nil
	}

	if err := createDirIfNotExist(m, cfg, stg); err != nil {
		return err
	}

	ctx := context.Background()

	lock, err := locker.Obtain(ctx, m.Name, LockingTimeTTL)
	if err != nil {
		return err
	}

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

func createDirIfNotExist(m mosaic.Mosaic, cfg *config.Config, stg storage.Storage) error {
	if cfg.StorageType.IsLocal() {
		return stg.CreateBucket(m.Name)
	}

	return nil
}
