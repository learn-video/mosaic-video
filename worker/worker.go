package worker

import (
	"context"
	"time"

	"github.com/mauricioabreu/mosaic-video/locking"
	"github.com/mauricioabreu/mosaic-video/mosaic"
)

func GenerateMosaic(key string, medias []mosaic.Media, locker locking.Locker, cmdExecutor mosaic.Command, runningProcesses map[string]bool) error {
	_, exists := runningProcesses[key]
	if exists {
		return nil
	}

	ctx := context.Background()
	lock, err := locker.Obtain(ctx, key, 5*time.Second)
	if err != nil {
		return err
	}

	cmdPath, args := mosaic.BuildCommand("ffmpeg", key, medias)
	if err := cmdExecutor.Execute(cmdPath, args...); err != nil {
		lock.Release(ctx)
		return err
	}

	runningProcesses[key] = true

	return nil
}
