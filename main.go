package main

import (
	"github.com/bsm/redislock"
	"github.com/mauricioabreu/mosaic-video/mosaic"
	"github.com/mauricioabreu/mosaic-video/worker"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugar := logger.Sugar()

	urls := []string{
		"https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/level_2.m3u8",
		"https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/level_4.m3u8",
	}

	rdc := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	locker := worker.NewRedisLocker(redislock.New(rdc))
	if err := worker.GenerateMosaic("test", urls, locker, &mosaic.FFMPEGCommand{}); err != nil {
		sugar.Fatal(err)
	}
}
