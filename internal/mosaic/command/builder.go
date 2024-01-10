package command

import (
	"fmt"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
)

func Build(mosaic mosaic.Mosaic, cfg *config.Config) []string {
	segmentPattern := fmt.Sprintf("%s/%s/seg_%%s.ts", cfg.AssetsPath, mosaic.Name)
	playlistPath := fmt.Sprintf("%s/%s/playlist.m3u8", cfg.AssetsPath, mosaic.Name)

	filterComplex := "nullsrc=size=1920x1080 [background];" +
		"[0:v] realtime, scale=1920x1080 [image];" +
		"[1:v] setpts=PTS-STARTPTS, scale=1170x660 [v1];" +
		"[2:v] setpts=PTS-STARTPTS, scale=568x320 [v2];" +
		"[background][v1] overlay=shortest=0:x=84:y=40 [posv1];" +
		"[posv1][v2] overlay=shortest=0:x=1260:y=40 [posv2];" +
		"[image][posv2] overlay=shortest=0 [mosaico]"

	args := []string{
		"-loglevel", "error",
		"-i", cfg.StaticsPath + "/background.jpg",
		"-i", mosaic.Medias[0].URL,
		"-i", mosaic.Medias[1].URL,
		"-filter_complex", filterComplex,
		"-map", "[mosaico]",
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
		"-hls_segment_filename", segmentPattern,
		playlistPath,
	}

	if mosaic.WithAudio {
		args = append(args, "-map", "1:a")
	}

	return args
}
