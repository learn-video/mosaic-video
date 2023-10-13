package worker

import (
	"context"
	"time"

	"github.com/mauricioabreu/mosaic-video/locking"
	"github.com/mauricioabreu/mosaic-video/mosaic"
)

func GenerateMosaic(key string, urls []string, locker locking.Locker, cmdExecutor mosaic.Command) error {
	ctx := context.Background()
	lock, err := locker.Obtain(ctx, key, 5*time.Second)
	if err != nil {
		return err
	}

	cmdPath, args := mosaic.BuildCommand("ffmpeg", key, urls)
	if err := cmdExecutor.Execute(cmdPath, args...); err != nil {
		lock.Release(ctx)
		return err
	}

	return nil
}
