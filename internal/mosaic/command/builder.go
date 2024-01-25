package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
)

func Build(m mosaic.Mosaic, cfg *config.Config) []string {
	videoInputArguments := getVideoInputArguments(m)
	audioInputArguments := getAudioInputArguments(m)
	filterComplexBuilder := getFilterComplexBuilder(m)
	filterComplexArguments := getFilterComplexArguments(filterComplexBuilder)
	encoderArguments := getEncoderArguments()
	audioGroupArguments := getAudioGroupArguments(m)

	var hlsArguments []string
	if cfg.StorageType == "local" {
		hlsArguments = getLocalSavingArguments(m, cfg)
	} else {
		hlsArguments = getHlsArguments(m, cfg)
	}

	args := []string{"-loglevel", "error"}
	args = append(args, videoInputArguments...)
	args = append(args, filterComplexArguments...)
	args = append(args, audioInputArguments...)
	args = append(args, audioGroupArguments...)
	args = append(args, encoderArguments...)
	args = append(args, hlsArguments...)

	return args
}

func getVideoInputArguments(m mosaic.Mosaic) []string {
	args := []string{"-i", m.BackgroundURL}

	for _, media := range m.Medias {
		if media.IsLoop {
			args = append(args,
				"-stream_loop", "-1",
			)
		}

		args = append(args,
			"-i", media.URL,
		)
	}

	return args
}

func getAudioInputArguments(m mosaic.Mosaic) []string {
	args := []string{}

	if m.Audio.IsNoAudio() {
		return args
	}

	if m.Audio.IsFirstInput() {
		args = append(args,
			"-map", "[a1] -c:a aac -b:a 128k",
		)
	}

	if m.Audio.IsAllInputs() {
		for i := range m.Medias {
			videoIndex := strconv.Itoa(i + 1)

			args = append(args,
				"-map", fmt.Sprintf("[a%s] -c:a aac -b:a 128k", videoIndex),
			)
		}
	}

	return args
}

func getFilterComplexBuilder(m mosaic.Mosaic) strings.Builder {
	var filter strings.Builder

	filter.WriteString("nullsrc=size=1920x1080 [background];")
	filter.WriteString("[0:v] realtime, scale=1920x1080 [image];")

	// Scale all videos
	for i, media := range m.Medias {
		videoIndex := strconv.Itoa(i + 1)

		// Scale each video and assign a label
		filter.WriteString(fmt.Sprintf("[%d:v] setpts=PTS-STARTPTS, scale=%s [v%s];", i+1, media.Scale, videoIndex))
	}

	// Then, overlay all videos
	lastOverlay := "[background]"

	for i := range m.Medias {
		videoIndex := strconv.Itoa(i + 1)

		x, y := m.Medias[i].Position.X, m.Medias[i].Position.Y

		filter.WriteString(fmt.Sprintf("%s[v%s] overlay=shortest=0:x=%d:y=%d [%s];", lastOverlay, videoIndex, x, y, "posv"+videoIndex))

		lastOverlay = "[posv" + videoIndex + "]"
	}

	filter.WriteString(fmt.Sprintf("[image]%s overlay=shortest=0 [mosaic]", lastOverlay))

	if !m.Audio.IsNoAudio() {
		filter.WriteString(";")

		for i := range m.Medias {
			videoIndex := strconv.Itoa(i + 1)

			filter.WriteString(fmt.Sprintf("[%s:a] aresample=async=1 [a%s]", videoIndex, videoIndex))

			if m.Audio.IsFirstInput() {
				break
			}

			if m.Audio.IsAllInputs() {
				if i < len(m.Medias)-1 {
					filter.WriteString(";")
				}
			}
		}
	}

	return filter
}

func getFilterComplexArguments(filter strings.Builder) []string {
	args := []string{
		"-filter_complex", filter.String(),
		"-map", "[mosaic]",
	}

	return args
}

func getAudioGroupArguments(m mosaic.Mosaic) []string {
	args := []string{}

	if m.Audio.IsNoAudio() {
		return args
	}

	if m.Audio.IsFirstInput() {
		args = append(args,
			"-var_stream_map",
			"a:0,agroup:audio,default:yes v:0,agroup:audio",
		)
	}

	if m.Audio.IsAllInputs() {
		group := strings.Builder{}

		for i := range m.Medias {
			if i == 0 {
				group.WriteString(fmt.Sprintf("a:%d,agroup:audio,default:yes ", i))
			} else {
				group.WriteString(fmt.Sprintf("a:%d,agroup:audio ", i))
			}
		}

		group.WriteString("v:0,agroup:audio")

		args = append(args,
			"-var_stream_map", group.String(),
		)
	}

	return args
}

func getEncoderArguments() []string {
	args := []string{
		"-c:v", "libx264",
		"-b:v", "1000k",
		"-x264opts", "keyint=30:min-keyint=30:scenecut=-1",
		"-preset", "ultrafast",
		"-threads", "0",
	}

	return args
}

func getHlsArguments(m mosaic.Mosaic, cfg *config.Config) []string {
	playlistPath := fmt.Sprintf("hls/%s/playlist_%%v.m3u8", m.Name)
	segmentPath := fmt.Sprintf("hls/%s/segment_%%v_%%03d.ts", m.Name)

	args := []string{
		"-f", "hls",
		"-hls_time", "5",
		"-hls_list_size", "6",
		"-sc_threshold", "0",
		"-method", "PUT",
		"-http_persistent", "1",
		"-hls_segment_filename", fmt.Sprintf("%s/%s", cfg.S3.UploaderEndpoint, segmentPath),
		"-master_pl_name", "master.m3u8",
		fmt.Sprintf("%s/%s", cfg.S3.UploaderEndpoint, playlistPath),
	}

	return args
}

func getLocalSavingArguments(m mosaic.Mosaic, cfg *config.Config) []string {
	playlistPath := fmt.Sprintf("%s/%s/playlist_%%v.m3u8", cfg.LocalStorage.Path, m.Name)
	segmentPath := fmt.Sprintf("%s/%s/segment_%%v_%%03d.ts", cfg.LocalStorage.Path, m.Name)

	args := []string{
		"-f", "hls",
		"-hls_time", "5",
		"-hls_list_size", "6",
		"-hls_start_number_source", "epoch",
		"-hls_segment_filename", segmentPath,
		"-master_pl_name", "master.m3u8",
		playlistPath,
	}

	return args
}
