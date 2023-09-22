package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	urls := []string{
		"https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/level_2.m3u8",
		"https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/level_4.m3u8",
	}
	cmd := buildCommand(urls)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error executing ffmpeg command: %v", err)
	}
}

func buildCommand(urls []string) *exec.Cmd {
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

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
