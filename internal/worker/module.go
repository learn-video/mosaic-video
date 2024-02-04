package worker

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hibiken/asynq"
	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/locking"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
	"github.com/mauricioabreu/mosaic-video/internal/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MosaicTaskParams struct {
	cfg    *config.Config
	logger *zap.SugaredLogger
	locker locking.Locker
	stg    storage.Storage
}

func Run(lc fx.Lifecycle, cfg *config.Config, logger *zap.SugaredLogger, locker *locking.RedisLocker, stg storage.Storage) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			redisAddress := cfg.Redis.Host + ":" + cfg.Redis.Port
			srv := asynq.NewServer(
				asynq.RedisClientOpt{Addr: redisAddress},
				asynq.Config{Concurrency: cfg.MaxConcurrentTasks},
			)

			mux := asynq.NewServeMux()

			startMosaicHandler := func(ctx context.Context, t *asynq.Task) error {
				return handleStartMosaicTask(ctx, t, MosaicTaskParams{
					cfg:    cfg,
					logger: logger,
					locker: locker,
					stg:    stg,
				})
			}
			mux.HandleFunc(TypeStartMosaic, startMosaicHandler)

			logger.Info("Worker started listening for tasks")

			go func() {
				if err := srv.Run(mux); err != nil {
					logger.Fatal("Could not run asynq server: ", err)
				}
			}()

			return nil
		},
	})
}

func handleStartMosaicTask(ctx context.Context, t *asynq.Task, mp MosaicTaskParams) error {
	var p StartMosaicPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	c := make(chan error, 1)
	go func() {
		mp.logger.Info("Processing mosaic: ", p.Mosaic.Name)

		c <- GenerateMosaic(
			ctx, p.Mosaic,
			mp.cfg,
			mp.logger,
			mp.locker,
			&mosaic.FFMPEGCommand{},
			mp.stg,
		)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c:
		if errors.Is(err, ErrLockFailed) {
			mp.logger.Debugf("Failed to obtain lock for mosaic '%s'", p.Mosaic.Name)
			return nil
		}

		mp.logger.Errorf("Error processing mosaic '%s': %w", p.Mosaic.Name, err)

		return err
	}
}
