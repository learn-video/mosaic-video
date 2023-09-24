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

func BuildCommand(commandPath string, key string, urls []string) (string, []string) {
	args := []string{}

	for _, url := range urls {
		args = append(args, "-i", url)
	}

	positions := []string{"0_0", "w0_0", "0_h0", "w0_h0"}
	xstackLayout := strings.Join(positions[:len(urls)], "|")

	inputLabels := make([]string, 0)
	filterComplex := ""
	for i := range urls {
		inputLabels = append(inputLabels, fmt.Sprintf("[l%d]", i))
		filterComplex += fmt.Sprintf("[%d:v] setpts=PTS-STARTPTS, scale=qvga %s; ", i, inputLabels[i])
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
