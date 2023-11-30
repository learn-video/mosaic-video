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

func Run(lc fx.Lifecycle, config *config.Config, logger *zap.SugaredLogger, locker *locking.RedisLocker) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				runningProcesses := make(map[string]string)
				for {
					logger.Info("worker is running")

					tasks, err := mosaic.FetchMosaicTasks(config.API.URL)
					if err != nil {
						logger.Fatal(err)
					}

					for _, task := range tasks {
						if err := GenerateMosaic(task.Name, task.Medias, locker, &mosaic.FFMPEGCommand{}, runningProcesses); err != nil {
							logger.Fatal(err)
						}
					}

					time.Sleep(120 * time.Second)
				}
			}()
			return nil
		},
	})
}
