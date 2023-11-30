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
		medias       []mosaic.Media
		expectedCmd  string
		expectedArgs []string
	}{
		{
			command:     "ffmpeg",
			name:        "Multiple URLs",
			key:         "mosaicvideo",
			medias:      []mosaic.Media{{URL: "http://example.com/mosaicvideo_1.m3u8", Position: "0_0"}, {URL: "http://example.com/mosaicvideo_2.m3u8", Position: "w0_0"}},
			expectedCmd: "ffmpeg",
			expectedArgs: []string{
				"-i", "http://example.com/mosaicvideo_1.m3u8",
				"-i", "http://example.com/mosaicvideo_2.m3u8",
				"-filter_complex", "[0:v] setpts=PTS-STARTPTS, scale=qvga [v0];[1:v] setpts=PTS-STARTPTS, scale=qvga [v1];[v0][v1]xstack=inputs=2:layout=0_0|w0_0[outv]; [0:a] aresample=async=1 [a0];[1:a] aresample=async=1 [a1]",
				"-map", "[outv]", "-c:v", "libx264", "-b:v", "800k", "-x264opts", "keyint=30:min-keyint=30:scenecut=-1",
				"-map", "[a0]", "-c:a", "aac", "-b:a", "128k",
				"-map", "[a1]", "-c:a", "aac", "-b:a", "128k",
				"-f", "hls",
				"-hls_time", "4",
				"-hls_list_size", "6",
				"-hls_flags", "delete_segments",
				"-hls_segment_filename", "output/seg_%v_%03d.ts",
				"-var_stream_map", "\"a:0,agroup:audio,language:ENG,default:yes a:1,agroup:audio,language:ENG v:0,agroup:audio\"",
				"-master_pl_name", "master.m3u8",
				"output/playlist_%v.m3u8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, args := mosaic.BuildCommand(tt.command, tt.key, tt.medias)
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
