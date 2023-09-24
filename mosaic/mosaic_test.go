package mosaic_test

import (
	"testing"

	"github.com/mauricioabreu/mosaic-video/mocks"
	"github.com/mauricioabreu/mosaic-video/mosaic"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBuildFFMPEGCommand(t *testing.T) {
	tests := []struct {
		command      string
		name         string
		key          string
		urls         []string
		expectedCmd  string
		expectedArgs []string
	}{
		{
			command:     "ffmpeg",
			name:        "Single URL",
			key:         "mosaicvideo",
			urls:        []string{"http://example.com/mosaicvideo.m3u8"},
			expectedCmd: "ffmpeg",
			expectedArgs: []string{
				"-i", "http://example.com/mosaicvideo.m3u8",
				"-map", "[out]",
				"-filter_complex", "[0:v] setpts=PTS-STARTPTS, scale=qvga [l0]; [l0]xstack=inputs=1:layout=0_0[out]",
				"-c:v", "libx264",
				"-x264opts", "keyint=30:min-keyint=30:scenecut=-1",
				"-f", "hls",
				"-hls_time", "5",
				"-hls_start_number_source", "epoch",
				"-hls_segment_filename", "output/segment%03d.ts",
				"output/playlist.m3u8",
			},
		},
		{
			command: "ffmpeg",
			name:    "Multiple URLs",
			key:     "mosaicvideo",
			urls: []string{
				"http://example.com/mosaicvideo_1.m3u8",
				"http://example.com/mosaicvideo_2.m3u8",
			},
			expectedCmd: "ffmpeg",
			expectedArgs: []string{
				"-i", "http://example.com/mosaicvideo_1.m3u8",
				"-i", "http://example.com/mosaicvideo_2.m3u8",
				"-map", "[out]",
				"-filter_complex", "[0:v] setpts=PTS-STARTPTS, scale=qvga [l0]; [1:v] setpts=PTS-STARTPTS, scale=qvga [l1]; [l0][l1]xstack=inputs=2:layout=0_0|w0_0[out]",
				"-c:v", "libx264",
				"-x264opts", "keyint=30:min-keyint=30:scenecut=-1",
				"-f", "hls",
				"-hls_time", "5",
				"-hls_start_number_source", "epoch",
				"-hls_segment_filename", "output/segment%03d.ts",
				"output/playlist.m3u8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, args := mosaic.BuildCommand(tt.command, tt.key, tt.urls)
			assert.Equal(t, tt.expectedCmd, cmd)
			assert.ElementsMatch(t, tt.expectedArgs, args)
		})
	}
}

func TestGenerateMosaic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmdExecutor := mocks.NewMockCommand(ctrl)
	cmdExecutor.EXPECT().Execute("ffmpeg", "arg1", "arg2").Return(nil)

	err := mosaic.GenerateMosaic(cmdExecutor, "ffmpeg", "arg1", "arg2")
	assert.NoError(t, err)
}
