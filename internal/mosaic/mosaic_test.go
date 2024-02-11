package mosaic_test

import (
	"context"
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
			name:    "Multiple URLs with first audio input and saving on cloud storage",
			key:     "mosaicvideo",
			mosaic: mosaic.Mosaic{
				Name:          "mosaicvideo",
				BackgroundURL: "http://example.com/background.jpg",
				Medias: []mosaic.Media{
					{
						URL: "http://example.com/mosaicvideo_1.m3u8",
						Position: mosaic.Position{
							X: 84,
							Y: 40,
						},
						Scale: "1170x660",
					},
					{
						URL: "http://example.com/mosaicvideo_2.m3u8",
						Position: mosaic.Position{
							X: 1260,
							Y: 40,
						},
						Scale: "568x320",
					},
				},
				Audio: mosaic.FirstInput,
			},
			cfg: &config.Config{
				StorageType: config.Cloud,
				S3:          config.S3{UploaderEndpoint: "http://localhost:8080"},
			},
			expectedCmd: "ffmpeg",
			expectedArgs: []string{
				"-loglevel", "error",
				"-i", "http://example.com/background.jpg",
				"-reconnect_at_eof", "1",
				"-reconnect_streamed", "1",
				"-reconnect_on_network_error", "1",
				"-reconnect_on_http_error", "4xx,5xx",
				"-reconnect_delay_max", "2",
				"-i", "http://example.com/mosaicvideo_1.m3u8",
				"-reconnect_at_eof", "1",
				"-reconnect_streamed", "1",
				"-reconnect_on_network_error", "1",
				"-reconnect_on_http_error", "4xx,5xx",
				"-reconnect_delay_max", "2",
				"-i", "http://example.com/mosaicvideo_2.m3u8",
				"-filter_complex", `nullsrc=size=1920x1080 [background];[0:v] realtime, scale=1920x1080 [image];[1:v] setpts=PTS-STARTPTS, scale=1170x660 [v1];[2:v] setpts=PTS-STARTPTS, scale=568x320 [v2];[background][v1] overlay=shortest=0:x=84:y=40 [posv1];[posv1][v2] overlay=shortest=0:x=1260:y=40 [posv2];[image][posv2] overlay=shortest=0 [mosaic];[1:a] aresample=async=1 [a1]`,
				"-map", "[mosaic]",
				"-map", "[a1] -c:a aac -b:a 128k",
				"-var_stream_map", "a:0,agroup:audio,default:yes v:0,agroup:audio",
				"-c:v", "libx264",
				"-b:v", "1000k",
				"-x264opts", "keyint=30:min-keyint=30:scenecut=-1",
				"-preset", "ultrafast",
				"-threads", "0",
				"-f", "hls",
				"-hls_time", "5",
				"-hls_list_size", "6",
				"-sc_threshold", "0",
				"-method", "PUT",
				"-http_persistent", "1",
				"-hls_segment_filename", "http://localhost:8080/hls/mosaicvideo/segment_%v_%03d.ts",
				"-master_pl_name", "master.m3u8",
				"http://localhost:8080/hls/mosaicvideo/playlist_%v.m3u8",
			},
		},
		{
			command: "ffmpeg",
			name:    "Multiple URLs with all audio inputs and saving on cloud storage",
			key:     "mosaicvideo",
			mosaic: mosaic.Mosaic{
				Name:          "mosaicvideo",
				BackgroundURL: "http://example.com/background.jpg",
				Medias: []mosaic.Media{
					{
						URL: "http://example.com/mosaicvideo_1.m3u8",
						Position: mosaic.Position{
							X: 84,
							Y: 40,
						},
						Scale: "1170x660",
					},
					{
						URL: "http://example.com/mosaicvideo_2.m3u8",
						Position: mosaic.Position{
							X: 1260,
							Y: 40,
						},
						Scale: "568x320",
					},
				},
				Audio: mosaic.AllInputs,
			},
			cfg: &config.Config{
				StorageType: config.Cloud,
				S3:          config.S3{UploaderEndpoint: "http://localhost:8080"},
			},
			expectedCmd: "ffmpeg",
			expectedArgs: []string{
				"-loglevel", "error",
				"-i", "http://example.com/background.jpg",
				"-reconnect_at_eof", "1",
				"-reconnect_streamed", "1",
				"-reconnect_on_network_error", "1",
				"-reconnect_on_http_error", "4xx,5xx",
				"-reconnect_delay_max", "2",
				"-i", "http://example.com/mosaicvideo_1.m3u8",
				"-reconnect_at_eof", "1",
				"-reconnect_streamed", "1",
				"-reconnect_on_network_error", "1",
				"-reconnect_on_http_error", "4xx,5xx",
				"-reconnect_delay_max", "2",
				"-i", "http://example.com/mosaicvideo_2.m3u8",
				"-filter_complex", `nullsrc=size=1920x1080 [background];[0:v] realtime, scale=1920x1080 [image];[1:v] setpts=PTS-STARTPTS, scale=1170x660 [v1];[2:v] setpts=PTS-STARTPTS, scale=568x320 [v2];[background][v1] overlay=shortest=0:x=84:y=40 [posv1];[posv1][v2] overlay=shortest=0:x=1260:y=40 [posv2];[image][posv2] overlay=shortest=0 [mosaic];[1:a] aresample=async=1 [a1];[2:a] aresample=async=1 [a2]`,
				"-map", "[mosaic]",
				"-map", "[a1] -c:a aac -b:a 128k",
				"-map", "[a2] -c:a aac -b:a 128k",
				"-var_stream_map", "a:0,agroup:audio,default:yes a:1,agroup:audio v:0,agroup:audio",
				"-c:v", "libx264",
				"-b:v", "1000k",
				"-x264opts", "keyint=30:min-keyint=30:scenecut=-1",
				"-preset", "ultrafast",
				"-threads", "0",
				"-f", "hls",
				"-hls_time", "5",
				"-hls_list_size", "6",
				"-sc_threshold", "0",
				"-method", "PUT",
				"-http_persistent", "1",
				"-hls_segment_filename", "http://localhost:8080/hls/mosaicvideo/segment_%v_%03d.ts",
				"-master_pl_name", "master.m3u8",
				"http://localhost:8080/hls/mosaicvideo/playlist_%v.m3u8",
			},
		},
		{
			command: "ffmpeg",
			name:    "Multiple URLs without audio and saving on cloud storage",
			key:     "mosaicvideo",
			mosaic: mosaic.Mosaic{
				Name:          "mosaicvideo",
				BackgroundURL: "http://example.com/background.jpg",
				Medias: []mosaic.Media{
					{
						URL: "http://example.com/mosaicvideo_1.m3u8",
						Position: mosaic.Position{
							X: 84,
							Y: 40,
						},
						Scale: "1170x660",
					},
					{
						URL: "http://example.com/mosaicvideo_2.m3u8",
						Position: mosaic.Position{
							X: 1260,
							Y: 40,
						},
						Scale: "568x320",
					},
				},
				Audio: mosaic.NoAudio,
			},
			cfg: &config.Config{
				StorageType: config.Cloud,
				S3:          config.S3{UploaderEndpoint: "http://localhost:8080"},
			},
			expectedCmd: "ffmpeg",
			expectedArgs: []string{
				"-loglevel", "error",
				"-i", "http://example.com/background.jpg",
				"-reconnect_at_eof", "1",
				"-reconnect_streamed", "1",
				"-reconnect_on_network_error", "1",
				"-reconnect_on_http_error", "4xx,5xx",
				"-reconnect_delay_max", "2",
				"-i", "http://example.com/mosaicvideo_1.m3u8",
				"-reconnect_at_eof", "1",
				"-reconnect_streamed", "1",
				"-reconnect_on_network_error", "1",
				"-reconnect_on_http_error", "4xx,5xx",
				"-reconnect_delay_max", "2",
				"-i", "http://example.com/mosaicvideo_2.m3u8",
				"-filter_complex", `nullsrc=size=1920x1080 [background];[0:v] realtime, scale=1920x1080 [image];[1:v] setpts=PTS-STARTPTS, scale=1170x660 [v1];[2:v] setpts=PTS-STARTPTS, scale=568x320 [v2];[background][v1] overlay=shortest=0:x=84:y=40 [posv1];[posv1][v2] overlay=shortest=0:x=1260:y=40 [posv2];[image][posv2] overlay=shortest=0 [mosaic]`,
				"-map", "[mosaic]",
				"-c:v", "libx264",
				"-b:v", "1000k",
				"-x264opts", "keyint=30:min-keyint=30:scenecut=-1",
				"-preset", "ultrafast",
				"-threads", "0",
				"-f", "hls",
				"-hls_time", "5",
				"-hls_list_size", "6",
				"-sc_threshold", "0",
				"-method", "PUT",
				"-http_persistent", "1",
				"-hls_segment_filename", "http://localhost:8080/hls/mosaicvideo/segment_%v_%03d.ts",
				"-master_pl_name", "master.m3u8",
				"http://localhost:8080/hls/mosaicvideo/playlist_%v.m3u8",
			},
		},
		{
			command: "ffmpeg",
			name:    "Multiple URLs with first audio input and saving in local storage",
			key:     "mosaicvideo",
			mosaic: mosaic.Mosaic{
				Name:          "mosaicvideo",
				BackgroundURL: "http://example.com/background.jpg",
				Medias: []mosaic.Media{
					{
						URL: "http://example.com/mosaicvideo_1.m3u8",
						Position: mosaic.Position{
							X: 84,
							Y: 40,
						},
						Scale: "1170x660",
					},
					{
						URL: "http://example.com/mosaicvideo_2.m3u8",
						Position: mosaic.Position{
							X: 1260,
							Y: 40,
						},
						Scale: "568x320",
					},
				},
				Audio: mosaic.FirstInput,
			},
			cfg: &config.Config{
				StorageType: config.Local,
				LocalStorage: config.LocalStorage{
					Path: "/home/hls",
				},
			},
			expectedCmd: "ffmpeg",
			expectedArgs: []string{
				"-loglevel", "error",
				"-i", "http://example.com/background.jpg",
				"-reconnect_at_eof", "1",
				"-reconnect_streamed", "1",
				"-reconnect_on_network_error", "1",
				"-reconnect_on_http_error", "4xx,5xx",
				"-reconnect_delay_max", "2",
				"-i", "http://example.com/mosaicvideo_1.m3u8",
				"-reconnect_at_eof", "1",
				"-reconnect_streamed", "1",
				"-reconnect_on_network_error", "1",
				"-reconnect_on_http_error", "4xx,5xx",
				"-reconnect_delay_max", "2",
				"-i", "http://example.com/mosaicvideo_2.m3u8",
				"-filter_complex", `nullsrc=size=1920x1080 [background];[0:v] realtime, scale=1920x1080 [image];[1:v] setpts=PTS-STARTPTS, scale=1170x660 [v1];[2:v] setpts=PTS-STARTPTS, scale=568x320 [v2];[background][v1] overlay=shortest=0:x=84:y=40 [posv1];[posv1][v2] overlay=shortest=0:x=1260:y=40 [posv2];[image][posv2] overlay=shortest=0 [mosaic];[1:a] aresample=async=1 [a1]`,
				"-map", "[mosaic]",
				"-map", "[a1] -c:a aac -b:a 128k",
				"-var_stream_map", "a:0,agroup:audio,default:yes v:0,agroup:audio",
				"-c:v", "libx264",
				"-b:v", "1000k",
				"-x264opts", "keyint=30:min-keyint=30:scenecut=-1",
				"-preset", "ultrafast",
				"-threads", "0",
				"-f", "hls",
				"-hls_time", "5",
				"-hls_list_size", "6",
				"-hls_start_number_source", "epoch",
				"-hls_segment_filename", "/home/hls/mosaicvideo/segment_%v_%03d.ts",
				"-master_pl_name", "master.m3u8",
				"/home/hls/mosaicvideo/playlist_%v.m3u8",
			},
		},
		{
			command: "ffmpeg",
			name:    "Multiple URLs with first audio input, saving in local storage and one VoD in looping",
			key:     "mosaicvideo",
			mosaic: mosaic.Mosaic{
				Name:          "mosaicvideo",
				BackgroundURL: "http://example.com/background.jpg",
				Medias: []mosaic.Media{
					{
						URL: "http://example.com/mosaicvideo_1.m3u8",
						Position: mosaic.Position{
							X: 84,
							Y: 40,
						},
						Scale: "1170x660",
					},
					{
						URL: "http://example.com/promotion-video.mp4",
						Position: mosaic.Position{
							X: 1260,
							Y: 40,
						},
						Scale:  "568x320",
						IsLoop: true,
					},
				},
				Audio: mosaic.FirstInput,
			},
			cfg: &config.Config{
				StorageType: config.Local,
				LocalStorage: config.LocalStorage{
					Path: "/home/hls",
				},
			},
			expectedCmd: "ffmpeg",
			expectedArgs: []string{
				"-loglevel", "error",
				"-i", "http://example.com/background.jpg",
				"-reconnect_at_eof", "1",
				"-reconnect_streamed", "1",
				"-reconnect_on_network_error", "1",
				"-reconnect_on_http_error", "4xx,5xx",
				"-reconnect_delay_max", "2",
				"-i", "http://example.com/mosaicvideo_1.m3u8",
				"-reconnect_at_eof", "1",
				"-reconnect_streamed", "1",
				"-reconnect_on_network_error", "1",
				"-reconnect_on_http_error", "4xx,5xx",
				"-reconnect_delay_max", "2",
				"-stream_loop", "-1",
				"-i", "http://example.com/promotion-video.mp4",
				"-filter_complex", `nullsrc=size=1920x1080 [background];[0:v] realtime, scale=1920x1080 [image];[1:v] setpts=PTS-STARTPTS, scale=1170x660 [v1];[2:v] setpts=PTS-STARTPTS, scale=568x320 [v2];[background][v1] overlay=shortest=0:x=84:y=40 [posv1];[posv1][v2] overlay=shortest=0:x=1260:y=40 [posv2];[image][posv2] overlay=shortest=0 [mosaic];[1:a] aresample=async=1 [a1]`,
				"-map", "[mosaic]",
				"-map", "[a1] -c:a aac -b:a 128k",
				"-var_stream_map", "a:0,agroup:audio,default:yes v:0,agroup:audio",
				"-c:v", "libx264",
				"-b:v", "1000k",
				"-x264opts", "keyint=30:min-keyint=30:scenecut=-1",
				"-preset", "ultrafast",
				"-threads", "0",
				"-f", "hls",
				"-hls_time", "5",
				"-hls_list_size", "6",
				"-hls_start_number_source", "epoch",
				"-hls_segment_filename", "/home/hls/mosaicvideo/segment_%v_%03d.ts",
				"-master_pl_name", "master.m3u8",
				"/home/hls/mosaicvideo/playlist_%v.m3u8",
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

	ctx := context.TODO()
	cmdExecutor := mocks.NewMockCommand(ctrl)
	cmdExecutor.EXPECT().Execute(ctx, "ffmpeg", "arg1", "arg2").Return(nil)

	err := mosaic.GenerateMosaic(ctx, cmdExecutor, "ffmpeg", "arg1", "arg2")
	assert.NoError(t, err)
}
