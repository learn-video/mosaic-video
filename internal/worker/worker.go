package worker

import (
	"context"
	"os"
	"time"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/locking"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic/command"
)

func createPath(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func GenerateMosaic(m mosaic.Mosaic, cfg *config.Config, locker locking.Locker, cmdExecutor mosaic.Command, runningProcesses map[string]bool) error {
	if err := createPath(cfg.AssetsPath + "/" + m.Name); err != nil {
		return err
	}

	_, exists := runningProcesses[m.Name]
	if exists {
		return nil
	}

	ctx := context.Background()
	lock, err := locker.Obtain(ctx, m.Name, 120*time.Second)
	if err != nil {
		return err
	}

	args := command.Build(m, cfg)
	if err := cmdExecutor.Execute("ffmpeg", args...); err != nil {
		if err := lock.Release(ctx); err != nil {
			return err
		}
		return err
	}

	runningProcesses[m.Name] = true

	return nil
}
