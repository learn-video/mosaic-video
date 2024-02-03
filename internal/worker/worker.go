package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/locking"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic/command"
	"github.com/mauricioabreu/mosaic-video/internal/storage"
	"go.uber.org/zap"
)

const (
	LockingTimeTTL    time.Duration = 120 * time.Second
	KeepAliveInterval time.Duration = LockingTimeTTL / 3
	ReleaseTimeout    time.Duration = 5 * time.Second
)

var (
	ErrLockFailed = fmt.Errorf("failed to obtain lock")
)

func GenerateMosaic(
	ctx context.Context,
	m mosaic.Mosaic,
	cfg *config.Config,
	logger *zap.SugaredLogger,
	locker locking.Locker,
	cmdExecutor mosaic.Command,
	runningProcesses *sync.Map,
	stg storage.Storage,
) error {
	lock, err := locker.Obtain(ctx, m.Name, LockingTimeTTL)
	if err != nil {
		return ErrLockFailed
	}

	defer func() {
		releaseCtx, cancel := context.WithTimeout(context.Background(), ReleaseTimeout)
		defer cancel()
		if releaseErr := lock.Release(releaseCtx); releaseErr != nil {
			logger.Errorf("Failed to release lock for %s: %v", m.Name, releaseErr)
		}
	}()

	if _, exists := runningProcesses.Load(m.Name); exists {
		return nil
	}

	runningProcesses.Store(m.Name, true)

	defer runningProcesses.Delete(m.Name)

	if err := createBucket(&m, cfg, stg); err != nil {
		return err
	}

	go keepAlive(ctx, logger, lock)

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		args := command.Build(m, cfg)
		if err := executeCommand(ctx, cmdExecutor, args); err != nil {
			if lerr := lock.Release(ctx); lerr != nil {
				return lerr
			}
		}
	}

	return nil
}

func executeCommand(ctx context.Context, cmdExecutor mosaic.Command, args []string) error {
	if err := cmdExecutor.Execute(ctx, "ffmpeg", args...); err != nil {
		return fmt.Errorf("failed to execute command, error=%w", err)
	}

	return nil
}

func createBucket(m *mosaic.Mosaic, cfg *config.Config, stg storage.Storage) error {
	if cfg.StorageType.IsLocal() {
		return stg.CreateBucket(m.Name)
	}

	return stg.CreateBucket(cfg.S3.BucketName)
}

func keepAlive(ctx context.Context, logger *zap.SugaredLogger, lock locking.Lock) {
	ticker := time.NewTicker(KeepAliveInterval)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			if err := lock.Refresh(ctx, LockingTimeTTL); err != nil {
				logger.Errorf("failed to refresh lock TTL, error=%w", err)
			}
		}
	}
}
