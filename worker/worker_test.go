package worker_test

import (
	"errors"
	"testing"

	"github.com/mauricioabreu/mosaic-video/mocks"
	"github.com/mauricioabreu/mosaic-video/mosaic"
	"github.com/mauricioabreu/mosaic-video/worker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGenerateMosaicWhenLockingFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	locker := mocks.NewMockLocker(ctrl)
	medias := []mosaic.Media{{URL: "http://mosaicvideos.com/video1.m3u8", Position: "0_0"}, {URL: "http://mosaicvideos.com/video2.m3u8", Position: "w0_0"}}
	locker.EXPECT().Obtain(gomock.Any(), "mosaicvideo1", gomock.Any()).Return(nil, errors.New("error obtaining lock"))
	runningProcesses := make(map[string]bool)

	err := worker.GenerateMosaic(
		"mosaicvideo1",
		medias,
		locker,
		nil,
		runningProcesses,
	)

	assert.Error(t, err)
}

func TestGenerateMosaicWhenExecutingCommandFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	locker := mocks.NewMockLocker(ctrl)
	medias := []mosaic.Media{{URL: "http://mosaicvideos.com/video1.m3u8", Position: "0_0"}, {URL: "http://mosaicvideos.com/video2.m3u8", Position: "w0_0"}}
	lock := mocks.NewMockLock(ctrl)
	lock.EXPECT().Release(gomock.Any()).Return(nil)
	locker.EXPECT().Obtain(gomock.Any(), "mosaicvideo1", gomock.Any()).Return(lock, nil)
	cmdExecutor := mocks.NewMockCommand(ctrl)
	cmdExecutor.EXPECT().Execute("ffmpeg", gomock.Any()).Return(errors.New("error executing command"))
	runningProcesses := make(map[string]bool)

	err := worker.GenerateMosaic(
		"mosaicvideo1",
		medias,
		locker,
		cmdExecutor,
		runningProcesses,
	)

	assert.Error(t, err)
}
