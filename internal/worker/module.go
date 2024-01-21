package worker

import (
	"context"
	"time"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/locking"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
	"github.com/mauricioabreu/mosaic-video/internal/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const SleepTime time.Duration = 60 * time.Second

func Run(lc fx.Lifecycle, cfg *config.Config, logger *zap.SugaredLogger, locker *locking.RedisLocker, stg storage.Storage) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				runningProcesses := make(map[string]bool)

				for {
					logger.Info("worker is running")

					tasks, err := mosaic.FetchMosaicTasks(cfg.API.URL)
					if err != nil {
						logger.Fatal(err)
					}

					for _, task := range tasks {
						go func(m mosaic.Mosaic) {
							defer func() {
								// Once finished, unmark the task
								delete(runningProcesses, m.Name)
							}()

							if err := GenerateMosaic(m, cfg, locker, &mosaic.FFMPEGCommand{}, runningProcesses, stg); err != nil {
								logger.Error(err)
							}
						}(task)
					}

					time.Sleep(SleepTime)
				}
			}()
			return nil
		},
	})
}
