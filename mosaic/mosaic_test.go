package mosaic_test

import (
	"testing"

	"github.com/mauricioabreu/mosaic-video/config"
	"github.com/mauricioabreu/mosaic-video/mocks"
	"github.com/mauricioabreu/mosaic-video/mosaic"
	"github.com/mauricioabreu/mosaic-video/mosaic/command"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBuildFFMPEGCommand(t *testing.T) {
	tests := []struct {
		command      string
		name         string
		key          string
		mosaic       mosaic.Mosaic
		cfg          *config.Config
		expectedCmd  string
		expectedArgs []string
	}{
		{
			command: "ffmpeg",
			name:    "Multiple URLs",
			key:     "mosaicvideo",
			mosaic: mosaic.Mosaic{
				Name: "mosaicvideo",
				Medias: []mosaic.Media{
					{URL: "http://example.com/mosaicvideo_1.m3u8"},
					{URL: "http://example.com/mosaicvideo_2.m3u8"},
				},
			},
			cfg: &config.Config{
				AssetsPath:  "output",
				StaticsPath: "statics",
			},
			expectedCmd: "ffmpeg",
			expectedArgs: []string{
				"-loglevel", "error",
				"-i", "statics/background.jpg",
				"-i", "http://example.com/mosaicvideo_1.m3u8",
				"-i", "http://example.com/mosaicvideo_2.m3u8",
				"-filter_complex", `nullsrc=size=1920x1080 [background];[0:v] realtime, scale=1920x1080 [image];[1:v] setpts=PTS-STARTPTS, scale=1170x660 [v1];[2:v] setpts=PTS-STARTPTS, scale=568x320 [v2];[background][v1] overlay=shortest=0:x=84:y=40 [posv1];[posv1][v2] overlay=shortest=0:x=1260:y=40 [posv2];[image][posv2] overlay=shortest=0 [mosaico]`,
				"-map", "[mosaico]",
				"-map", "1:a",
				"-c:v", "libx264",
				"-x264opts", "keyint=30:min-keyint=30:scenecut=-1",
				"-preset", "ultrafast",
				"-threads", "0",
				"-r", "24",
				"-c:a", "copy",
				"-f", "hls",
				"-hls_time", "5",
				"-hls_list_size", "12",
				"-hls_flags", "delete_segments",
				"-strftime", "1",
				"-hls_segment_filename", "output/mosaicvideo/seg_%s.ts",
				"output/mosaicvideo/playlist.m3u8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := command.Build(tt.mosaic, tt.cfg)
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
