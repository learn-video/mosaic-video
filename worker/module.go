package worker

import (
	"context"
	"time"

	"github.com/mauricioabreu/mosaic-video/locking"
	"github.com/mauricioabreu/mosaic-video/mosaic"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(lc fx.Lifecycle, logger *zap.SugaredLogger, locker *locking.RedisLocker) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				for {
					logger.Info("worker is running")
					urls := []string{
						"https://ireplay.tv/test/rate_5_28.m3u8",
						"https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/level_4.m3u8",
					}

					if err := GenerateMosaic("test", urls, locker, &mosaic.FFMPEGCommand{}); err != nil {
						logger.Fatal(err)
					}
					time.Sleep(120 * time.Second)
				}
			}()
			return nil
		},
	})
}
