package worker

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/mauricioabreu/mosaic-video/mosaic"
)

type Locker interface {
	Obtain(ctx context.Context, key string, ttl time.Duration) (Lock, error)
}

type Lock interface {
	Release(ctx context.Context) error
	TTL(ctx context.Context) (time.Duration, error)
	Refresh(ctx context.Context, ttl time.Duration) error
}

type RedisLocker struct {
	locker *redislock.Client
}

func NewRedisLocker(client *redislock.Client) *RedisLocker {
	return &RedisLocker{locker: client}
}

func (r *RedisLocker) Obtain(ctx context.Context, key string, ttl time.Duration) (Lock, error) {
	lock, err := r.locker.Obtain(ctx, key, ttl, nil)
	if err != nil {
		return nil, err
	}
	return &RedisLock{lock: lock}, nil
}

type RedisLock struct {
	lock *redislock.Lock
}

func (r *RedisLock) Release(ctx context.Context) error {
	return r.lock.Release(ctx)
}

func (r *RedisLock) TTL(ctx context.Context) (time.Duration, error) {
	return r.lock.TTL(ctx)
}

func (r *RedisLock) Refresh(ctx context.Context, ttl time.Duration) error {
	return r.lock.Refresh(ctx, ttl, nil)
}

func GenerateMosaic(key string, urls []string, locker Locker, cmdExecutor mosaic.Command) error {
	ctx := context.Background()
	lock, err := locker.Obtain(ctx, key, 5*time.Second)
	if err != nil {
		return err
	}

	cmdPath, args := mosaic.BuildCommand("ffmpeg", key, urls)
	if err := cmdExecutor.Execute(cmdPath, args...); err != nil {
		lock.Release(ctx)
		return err
	}

	return nil
}
