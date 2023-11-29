package mosaic

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Command interface {
	Execute(command string, args ...string) error
}

type FFMPEGCommand struct{}

func (r *FFMPEGCommand) Execute(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func GenerateMosaic(executor Command, command string, args ...string) error {
	return executor.Execute(command, args...)
}

func BuildCommand(commandPath string, key string, medias []Media) (string, []string) {
	args := []string{}

	for _, media := range medias {
		args = append(args, "-i", media.URL)
	}

	positions := make([]string, len(medias))
	for i, media := range medias {
		positions[i] = media.Position
	}
	xstackLayout := strings.Join(positions, "|")

	inputLabels := make([]string, 0)
	filterComplex := ""
	for i := range medias {
		inputLabels = append(inputLabels, fmt.Sprintf("[l%d]", i))
		filterComplex += fmt.Sprintf("[%d:v] setpts=PTS-STARTPTS, scale=qvga %s; ", i, inputLabels[i])
	}

	urls := make([]string, 0)
	for _, url := range medias {
		urls = append(urls, url.URL)
	}

	filterComplex += fmt.Sprintf("%sxstack=inputs=%d:layout=%s[out]", strings.Join(inputLabels, ""), len(urls), xstackLayout)

	args = append(args,
		"-filter_complex", filterComplex,
		"-map", "[out]",
		"-c:v", "libx264",
		"-x264opts", "keyint=30:min-keyint=30:scenecut=-1",
		"-f", "hls",
		"-hls_time", "5",
		"-hls_start_number_source", "epoch",
		"-hls_segment_filename", "output/segment%03d.ts",
		"output/playlist.m3u8",
	)

	return commandPath, args
}
