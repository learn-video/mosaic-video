package worker_test

import (
	"errors"
	"testing"

	"github.com/mauricioabreu/mosaic-video/mocks"
	"github.com/mauricioabreu/mosaic-video/worker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGenerateMosaicWhenLockingFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	locker := mocks.NewMockLocker(ctrl)
	urls := []string{"http://mosaicvideos.com/video1.m3u8", "http://mosaicvideos.com/video2.m3u8"}
	locker.EXPECT().Obtain(gomock.Any(), "mosaicvideo1", gomock.Any()).Return(nil, errors.New("error obtaining lock"))
	runningProcesses := make(map[string]string)

	err := worker.GenerateMosaic(
		"mosaicvideo1",
		urls,
		locker,
		nil,
		runningProcesses,
	)

	assert.Error(t, err)
}

func TestGenerateMosaicWhenExecutingCommandFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	locker := mocks.NewMockLocker(ctrl)
	urls := []string{"http://mosaicvideos.com/video1.m3u8", "http://mosaicvideos.com/video2.m3u8"}
	lock := mocks.NewMockLock(ctrl)
	lock.EXPECT().Release(gomock.Any()).Return(nil)
	locker.EXPECT().Obtain(gomock.Any(), "mosaicvideo1", gomock.Any()).Return(lock, nil)
	cmdExecutor := mocks.NewMockCommand(ctrl)
	cmdExecutor.EXPECT().Execute("ffmpeg", gomock.Any()).Return(errors.New("error executing command"))
	runningProcesses := make(map[string]string)

	err := worker.GenerateMosaic(
		"mosaicvideo1",
		urls,
		locker,
		cmdExecutor,
		runningProcesses,
	)

	assert.Error(t, err)
}
