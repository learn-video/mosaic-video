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
	defer ctrl.Finish()

	locker := mocks.NewMockLocker(ctrl)
	storage := mocks.NewMockStorage(ctrl)
	mosaic := mosaic.Mosaic{
		Name: "mosaicvideo",
		Medias: []mosaic.Media{
			{URL: "http://example.com/mosaicvideo_1.m3u8"},
			{URL: "http://example.com/mosaicvideo_2.m3u8"},
		},
	}
	cfg := &config.Config{}
	locker.EXPECT().Obtain(gomock.Any(), "mosaicvideo", gomock.Any()).Return(nil, errors.New("error obtaining lock"))
	storage.EXPECT().CreateBucket(gomock.Any()).Return(nil)

	runningProcesses := make(map[string]bool)

	err := worker.GenerateMosaic(
		mosaic,
		cfg,
		locker,
		nil,
		runningProcesses,
		storage,
	)

	assert.Error(t, err)
}

func TestGenerateMosaicWhenExecutingCommandFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	locker := mocks.NewMockLocker(ctrl)
	storage := mocks.NewMockStorage(ctrl)
	mosaic := mosaic.Mosaic{
		Name: "mosaicvideo",
		Medias: []mosaic.Media{
			{
				URL: "http://example.com/mosaicvideo_1.m3u8",
				Position: mosaic.Position{
					X: 84,
					Y: 40,
				},
				Scale: "1170x660"},
			{
				URL: "http://example.com/mosaicvideo_2.m3u8",
				Position: mosaic.Position{
					X: 1260,
					Y: 40,
				},
				Scale: "568x320",
			},
		},
	}
	lock := mocks.NewMockLock(ctrl)
	lock.EXPECT().Release(gomock.Any()).Return(nil)
	locker.EXPECT().Obtain(gomock.Any(), "mosaicvideo", gomock.Any()).Return(lock, nil)
	storage.EXPECT().CreateBucket(gomock.Any()).Return(nil)

	cmdExecutor := mocks.NewMockCommand(ctrl)
	cmdExecutor.EXPECT().Execute("ffmpeg", gomock.Any()).Return(errors.New("error executing command"))
	runningProcesses := make(map[string]bool)

	err := worker.GenerateMosaic(
		mosaic,
		cfg,
		locker,
		cmdExecutor,
		runningProcesses,
		storage,
	)

	assert.Error(t, err)
}

func TestGenerateMosaicWhenCreateBucketFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	storage := mocks.NewMockStorage(ctrl)
	mosaic := mosaic.Mosaic{
		Name: "mosaicvideo",
		Medias: []mosaic.Media{
			{
				URL: "http://example.com/mosaicvideo_1.m3u8",
				Position: mosaic.Position{
					X: 84,
					Y: 40,
				},
				Scale: "1170x660"},
			{
				URL: "http://example.com/mosaicvideo_2.m3u8",
				Position: mosaic.Position{
					X: 1260,
					Y: 40,
				},
				Scale: "568x320",
			},
		},
	}

	storage.EXPECT().CreateBucket(gomock.Any()).Return(errors.New("no permissions to create directory"))
	runningProcesses := make(map[string]bool)

	err := worker.GenerateMosaic(
		mosaic,
		cfg,
		nil,
		nil,
		runningProcesses,
		storage,
	)

	assert.Error(t, err)
}
