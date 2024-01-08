package worker_test

import (
	"errors"
	"testing"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/mocks"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
	"github.com/mauricioabreu/mosaic-video/internal/worker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGenerateMosaicWhenLockingFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	locker := mocks.NewMockLocker(ctrl)
	watcher := mocks.NewMockWatcher(ctrl)
	mosaic := mosaic.Mosaic{
		Name: "mosaicvideo",
		Medias: []mosaic.Media{
			{URL: "http://example.com/mosaicvideo_1.m3u8"},
			{URL: "http://example.com/mosaicvideo_2.m3u8"},
		},
	}
	cfg := &config.Config{AssetsPath: "output"}
	locker.EXPECT().Obtain(gomock.Any(), "mosaicvideo", gomock.Any()).Return(nil, errors.New("error obtaining lock"))
	runningProcesses := make(map[string]bool)

	err := worker.GenerateMosaic(
		mosaic,
		cfg,
		locker,
		nil,
		runningProcesses,
		watcher,
	)

	assert.Error(t, err)
}

func TestGenerateMosaicWhenExecutingCommandFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	locker := mocks.NewMockLocker(ctrl)
	watcher := mocks.NewMockWatcher(ctrl)
	mosaic := mosaic.Mosaic{
		Name: "mosaicvideo",
		Medias: []mosaic.Media{
			{URL: "http://example.com/mosaicvideo_1.m3u8"},
			{URL: "http://example.com/mosaicvideo_2.m3u8"},
		},
	}
	cfg := &config.Config{AssetsPath: "output"}
	lock := mocks.NewMockLock(ctrl)
	lock.EXPECT().Release(gomock.Any()).Return(nil)
	locker.EXPECT().Obtain(gomock.Any(), "mosaicvideo", gomock.Any()).Return(lock, nil)

	cmdExecutor := mocks.NewMockCommand(ctrl)
	cmdExecutor.EXPECT().Execute("ffmpeg", gomock.Any()).Return(errors.New("error executing command"))
	runningProcesses := make(map[string]bool)

	err := worker.GenerateMosaic(
		mosaic,
		cfg,
		locker,
		cmdExecutor,
		runningProcesses,
		watcher,
	)

	assert.Error(t, err)
}
