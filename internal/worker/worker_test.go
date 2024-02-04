package worker_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/logging"
	"github.com/mauricioabreu/mosaic-video/internal/mocks"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
	"github.com/mauricioabreu/mosaic-video/internal/worker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGenerateMosaicSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	logger := logging.NewLogger()
	locker := mocks.NewMockLocker(ctrl)
	storage := mocks.NewMockStorage(ctrl)
	cmdExecutor := mocks.NewMockCommand(ctrl)
	mosaic := mosaic.Mosaic{
		Name: "mosaicvideo",
		Medias: []mosaic.Media{
			{URL: "http://example.com/mosaicvideo_1.m3u8"},
			{URL: "http://example.com/mosaicvideo_2.m3u8"},
		},
	}

	lock := mocks.NewMockLock(ctrl)
	lock.EXPECT().Release(gomock.Any()).AnyTimes().Return(nil)
	locker.EXPECT().Obtain(gomock.Any(), mosaic.Name, gomock.Any()).Return(lock, nil)
	storage.EXPECT().CreateBucket(gomock.Any()).Return(nil)
	cmdExecutor.EXPECT().Execute(gomock.Any(), "ffmpeg", gomock.Any()).Return(nil)

	runningProcesses := &sync.Map{}

	err := worker.GenerateMosaic(
		context.TODO(),
		mosaic,
		cfg,
		logger,
		locker,
		cmdExecutor,
		runningProcesses,
		storage,
	)

	assert.NoError(t, err)
}

func TestGenerateMosaicWhenLockingFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	logger := logging.NewLogger()
	locker := mocks.NewMockLocker(ctrl)
	storage := mocks.NewMockStorage(ctrl)
	mosaic := mosaic.Mosaic{
		Name: "mosaicvideo",
		Medias: []mosaic.Media{
			{URL: "http://example.com/mosaicvideo_1.m3u8"},
			{URL: "http://example.com/mosaicvideo_2.m3u8"},
		},
	}
	locker.EXPECT().Obtain(gomock.Any(), "mosaicvideo", gomock.Any()).Return(nil, errors.New("error obtaining lock"))

	runningProcesses := &sync.Map{}

	err := worker.GenerateMosaic(
		context.TODO(),
		mosaic,
		cfg,
		logger,
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
	logger := logging.NewLogger()
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
	lock.EXPECT().Release(gomock.Any()).Return(nil)

	ctx := context.TODO()
	cmdExecutor := mocks.NewMockCommand(ctrl)
	cmdExecutor.EXPECT().Execute(ctx, "ffmpeg", gomock.Any()).Return(errors.New("error executing command"))

	runningProcesses := &sync.Map{}

	err := worker.GenerateMosaic(
		ctx,
		mosaic,
		cfg,
		logger,
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
	logger := logging.NewLogger()
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
	storage.EXPECT().CreateBucket(gomock.Any()).Return(errors.New("no permissions to create directory"))

	runningProcesses := &sync.Map{}

	err := worker.GenerateMosaic(
		context.TODO(),
		mosaic,
		cfg,
		logger,
		locker,
		nil,
		runningProcesses,
		storage,
	)

	assert.Error(t, err)
}
