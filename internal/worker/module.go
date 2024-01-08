package worker

import (
	"context"
	"time"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/locking"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
	"github.com/mauricioabreu/mosaic-video/internal/watcher"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(lc fx.Lifecycle, cfg *config.Config, logger *zap.SugaredLogger, locker *locking.RedisLocker, w watcher.Watcher) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				runningProcesses := make(map[string]bool)

				w.Run()

				go func() {
					for {
						select {
						case event := <-w.Events():
							logger.Infof("File system event: %v", event)
						case err := <-w.Errors():
							logger.Errorf("File system error: %v", err)
						}
					}
				}()

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

							if err := GenerateMosaic(m, cfg, locker, &mosaic.FFMPEGCommand{}, runningProcesses, w); err != nil {
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
