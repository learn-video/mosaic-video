package worker

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/locking"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
	"github.com/mauricioabreu/mosaic-video/internal/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(lc fx.Lifecycle, cfg *config.Config, logger *zap.SugaredLogger, locker *locking.RedisLocker, stg storage.Storage) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			redisAddress := cfg.Redis.Host + ":" + cfg.Redis.Port
			srv := asynq.NewServer(
				asynq.RedisClientOpt{Addr: redisAddress},
				asynq.Config{Concurrency: cfg.MaxConcurrentTasks},
			)

			rp := make(map[string]bool)

			mux := asynq.NewServeMux()

			startMosaicHandler := func(ctx context.Context, t *asynq.Task) error {
				return handleStartMosaicTask(ctx, t, cfg, logger, locker, rp, stg)
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

func handleStartMosaicTask(ctx context.Context, t *asynq.Task, cfg *config.Config, logger *zap.SugaredLogger, locker locking.Locker, rp map[string]bool, stg storage.Storage) error {
	var p StartMosaicPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	c := make(chan error, 1)
	go func() {
		logger.Info("Processing mosaic: ", p.Mosaic.Name)

		c <- GenerateMosaic(ctx, p.Mosaic, cfg, logger, locker, &mosaic.FFMPEGCommand{}, rp, stg)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c:
		return err
	}
}
