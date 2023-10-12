package main

import (
	"github.com/bsm/redislock"
	"github.com/mauricioabreu/mosaic-video/mosaic"
	"github.com/mauricioabreu/mosaic-video/worker"
	"github.com/redis/go-redis/v9"
)

func main() {
	urls := []string{
		"https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/level_2.m3u8",
		"https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/level_4.m3u8",
	}

	rdc := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	locker := worker.NewRedisLocker(redislock.New(rdc))
	worker.GenerateMosaic("test", urls, locker, &mosaic.FFMPEGCommand{})
}
