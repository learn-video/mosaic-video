package worker

import (
	"context"
	"os"
	"time"

	"github.com/mauricioabreu/mosaic-video/config"
	"github.com/mauricioabreu/mosaic-video/locking"
	"github.com/mauricioabreu/mosaic-video/mosaic"
	"github.com/mauricioabreu/mosaic-video/mosaic/command"
)

func createPath(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func GenerateMosaic(mosaic mosaic.Mosaic, cfg *config.Config, locker locking.Locker, cmdExecutor mosaic.Command, runningProcesses map[string]bool) error {
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
