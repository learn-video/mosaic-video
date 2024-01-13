package mosaic_test

import (
	"testing"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/mocks"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic/command"
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
			name:    "Multiple URLs with audio",
			key:     "mosaicvideo",
			mosaic: mosaic.Mosaic{
				Name: "mosaicvideo",
				Medias: []mosaic.Media{
					{URL: "http://example.com/mosaicvideo_1.m3u8", Position: "84_40", Scale: "1170x660"},
					{URL: "http://example.com/mosaicvideo_2.m3u8", Position: "1260_40", Scale: "568x320"},
				},
				WithAudio: true,
			},
			cfg: &config.Config{
				StaticsPath:      "statics",
				UploaderEndpoint: "http://localhost:8080",
			},
			expectedCmd: "ffmpeg",
			expectedArgs: []string{
				"-loglevel", "error",
				"-i", "statics/background.jpg",
				"-i", "http://example.com/mosaicvideo_1.m3u8",
				"-i", "http://example.com/mosaicvideo_2.m3u8",
				"-filter_complex", `nullsrc=size=1920x1080 [background];[0:v] realtime, scale=1920x1080 [image];[1:v] setpts=PTS-STARTPTS, scale=1170x660 [v1];[2:v] setpts=PTS-STARTPTS, scale=568x320 [v2];[background][v1] overlay=shortest=0:x=84:y=40 [posv1];[posv1][v2] overlay=shortest=0:x=1260:y=40 [posv2];[image][posv2] overlay=shortest=0 [mosaic]`,
				"-map", "[mosaic]",
				"-map", "1:a",
				"-c:v", "libx264",
				"-x264opts", "keyint=30:min-keyint=30:scenecut=-1",
				"-preset", "ultrafast",
				"-threads", "0",
				"-r", "24",
				"-c:a", "copy",
				"-f", "hls",
				"-hls_playlist_type", "event",
				"-hls_time", "5",
				"-strftime", "1",
				"-method", "PUT",
				"-http_persistent", "1",
				"-sc_threshold", "0",
				"http://localhost:8080/hls/mosaicvideo/playlist.m3u8",
			},
		},
		{
			command: "ffmpeg",
			name:    "Multiple URLs without audio",
			key:     "mosaicvideo",
			mosaic: mosaic.Mosaic{
				Name: "mosaicvideo",
				Medias: []mosaic.Media{
					{URL: "http://example.com/mosaicvideo_1.m3u8", Position: "84_40", Scale: "1170x660"},
					{URL: "http://example.com/mosaicvideo_2.m3u8", Position: "1260_40", Scale: "568x320"},
				},
				WithAudio: false,
			},
			cfg: &config.Config{
				StaticsPath:      "statics",
				UploaderEndpoint: "http://localhost:8080",
			},
			expectedCmd: "ffmpeg",
			expectedArgs: []string{
				"-loglevel", "error",
				"-i", "statics/background.jpg",
				"-i", "http://example.com/mosaicvideo_1.m3u8",
				"-i", "http://example.com/mosaicvideo_2.m3u8",
				"-filter_complex", `nullsrc=size=1920x1080 [background];[0:v] realtime, scale=1920x1080 [image];[1:v] setpts=PTS-STARTPTS, scale=1170x660 [v1];[2:v] setpts=PTS-STARTPTS, scale=568x320 [v2];[background][v1] overlay=shortest=0:x=84:y=40 [posv1];[posv1][v2] overlay=shortest=0:x=1260:y=40 [posv2];[image][posv2] overlay=shortest=0 [mosaic]`,
				"-map", "[mosaic]",
				"-c:v", "libx264",
				"-x264opts", "keyint=30:min-keyint=30:scenecut=-1",
				"-preset", "ultrafast",
				"-threads", "0",
				"-r", "24",
				"-c:a", "copy",
				"-f", "hls",
				"-hls_playlist_type", "event",
				"-hls_time", "5",
				"-strftime", "1",
				"-method", "PUT",
				"-http_persistent", "1",
				"-sc_threshold", "0",
				"http://localhost:8080/hls/mosaicvideo/playlist.m3u8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := command.Build(tt.mosaic, tt.cfg)
			assert.Equal(t, tt.expectedArgs, args)
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
