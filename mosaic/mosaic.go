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

	videoFilters := make([]string, 0)
	audioFilters := make([]string, 0)
	for i := range medias {
		videoLabel := fmt.Sprintf("[v%d]", i)
		audioLabel := fmt.Sprintf("[a%d]", i)
		videoFilters = append(videoFilters, fmt.Sprintf("[%d:v] setpts=PTS-STARTPTS, scale=qvga %s;", i, videoLabel))
		audioFilters = append(audioFilters, fmt.Sprintf("[%d:a] aresample=async=1 %s", i, audioLabel))
	}

	positions := make([]string, len(medias))
	for i, media := range medias {
		positions[i] = media.Position
	}

	xstackInputs := make([]string, len(medias))
	for i := range medias {
		xstackInputs[i] = fmt.Sprintf("[v%d]", i)
	}
	xstackLayout := strings.Join(positions, "|")
	filterComplex := fmt.Sprintf("%s%sxstack=inputs=%d:layout=%s[outv]; %s",
		strings.Join(videoFilters, ""),
		strings.Join(xstackInputs, ""),
		len(medias),
		xstackLayout,
		strings.Join(audioFilters, ";"))

	args = append(args,
		"-filter_complex", filterComplex,
		"-map", "[outv]", "-c:v", "libx264", "-b:v", "800k", "-x264opts", "keyint=30:min-keyint=30:scenecut=-1",
	)

	for i := range medias {
		args = append(args, "-map", fmt.Sprintf("[a%d]", i), "-c:a", "aac", "-b:a", "128k")
	}

	varStreamMapParts := make([]string, 0)
	for i := range medias {
		varStreamMapPart := fmt.Sprintf("a:%d,agroup:audio,language:ENG", i)
		if i == 0 {
			varStreamMapPart += ",default:yes"
		}
		varStreamMapParts = append(varStreamMapParts, varStreamMapPart)
	}
	varStreamMapParts = append(varStreamMapParts, "v:0,agroup:audio")
	varStreamMap := strings.Join(varStreamMapParts, " ")

	args = append(args,
		"-f", "hls",
		"-hls_time", "4",
		"-hls_list_size", "6",
		"-hls_flags", "delete_segments",
		"-hls_segment_filename", "output/seg_%v_%03d.ts",
		"-var_stream_map", varStreamMap,
		"-master_pl_name", "master.m3u8",
		"output/playlist_%v.m3u8",
	)

	return commandPath, args
}
