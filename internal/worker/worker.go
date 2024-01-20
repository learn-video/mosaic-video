package worker

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/locking"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic/command"
)

const LockingTimeTTL time.Duration = 120 * time.Second

func GenerateMosaic(m mosaic.Mosaic, cfg *config.Config, locker locking.Locker, cmdExecutor mosaic.Command, runningProcesses map[string]bool) error {
	_, exists := runningProcesses[m.Name]
	if exists {
		return nil
	}

	if err := createDirIfNotExist(m, cfg); err != nil {
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

func createDirIfNotExist(m mosaic.Mosaic, cfg *config.Config) error {
	if cfg.StorageType == "local" {
		path := fmt.Sprintf("%s/%s", cfg.LocalStorage.Path, m.Name)

		if _, err := os.Stat(path); err != nil {
			if !os.IsNotExist(err) {
				return err
			}

			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create mosaic directory path=%s  error=%w", path, err)
			}
		}
	}

	return nil
}
