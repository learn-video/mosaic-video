package worker

import (
	"context"
	"time"

	"github.com/mauricioabreu/mosaic-video/config"
	"github.com/mauricioabreu/mosaic-video/locking"
	"github.com/mauricioabreu/mosaic-video/mosaic"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(lc fx.Lifecycle, cfg *config.Config, logger *zap.SugaredLogger, locker *locking.RedisLocker) {
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
								delete(runningProcesses, task.Name)
							}()

							if err := GenerateMosaic(m, cfg, locker, &mosaic.FFMPEGCommand{}, runningProcesses); err != nil {
								logger.Error(err)
							}
						}(task)
					}

					time.Sleep(60 * time.Second)
				}
			}()
			return nil
		},
	})
}
