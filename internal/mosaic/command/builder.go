package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
)

func Build(m mosaic.Mosaic, cfg *config.Config) []string {
	playlistPath := fmt.Sprintf("hls/%s/playlist.m3u8", m.Name)

	var filterComplexBuilder strings.Builder
	filterComplexBuilder.WriteString("nullsrc=size=1920x1080 [background];")
	filterComplexBuilder.WriteString("[0:v] realtime, scale=1920x1080 [image];")

	// Scale all videos
	for i, media := range m.Medias {
		videoIndex := strconv.Itoa(i + 1)

		// Scale each video and assign a label
		filterComplexBuilder.WriteString(fmt.Sprintf("[%d:v] setpts=PTS-STARTPTS, scale=%s [v%s];", i+1, media.Scale, videoIndex))
	}

	// Then, overlay all videos
	lastOverlay := "[background]"
	for i := range m.Medias {
		videoIndex := strconv.Itoa(i + 1)

		x, y := m.Medias[i].Position.X, m.Medias[i].Position.Y

		filterComplexBuilder.WriteString(fmt.Sprintf("%s[v%s] overlay=shortest=0:x=%d:y=%d [%s];", lastOverlay, videoIndex, x, y, "posv"+videoIndex))

		lastOverlay = "[posv" + videoIndex + "]"
	}

	filterComplexBuilder.WriteString(fmt.Sprintf("[image]%s overlay=shortest=0 [mosaic]", lastOverlay))

	args := []string{
		"-loglevel", "error",
		"-i", m.BackgroundURL,
		"-i", m.Medias[0].URL,
		"-i", m.Medias[1].URL,
		"-filter_complex", filterComplexBuilder.String(),
		"-map", "[mosaic]",
	}

	if m.WithAudio {
		args = append(args, "-map", "1:a")
	}

	args = append(args, []string{
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
		fmt.Sprintf("%s/%s", cfg.UploaderEndpoint, playlistPath),
	}...)

	return args
}
