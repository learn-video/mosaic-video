package worker

import (
	"context"
	"os"
	"time"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/locking"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic/command"
	"github.com/mauricioabreu/mosaic-video/internal/watcher"
)

func createPath(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func GenerateMosaic(mosaic mosaic.Mosaic, cfg *config.Config, locker locking.Locker, cmdExecutor mosaic.Command, runningProcesses map[string]bool, w watcher.Watcher) error {
	if err := createPath(cfg.AssetsPath + "/" + mosaic.Name); err != nil {
		return err
	}

	_, exists := runningProcesses[mosaic.Name]
	if exists {
		return nil
	}

	ctx := context.Background()
	lock, err := locker.Obtain(ctx, mosaic.Name, 120*time.Second)
	if err != nil {
		return err
	}

	args := command.Build(mosaic, cfg)
	if err := cmdExecutor.Execute("ffmpeg", args...); err != nil {
		lock.Release(ctx)
		return err
	}

	runningProcesses[mosaic.Name] = true

	return nil
}
