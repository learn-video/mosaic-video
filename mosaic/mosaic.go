package mosaic

import (
	"os"
	"os/exec"
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
