package mosaic

import (
	"context"
	"os"
	"os/exec"
)

type Command interface {
	Execute(ctx context.Context, command string, args ...string) error
}

type FFMPEGCommand struct{}

func (r *FFMPEGCommand) Execute(ctx context.Context, command string, args ...string) error {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

func GenerateMosaic(ctx context.Context, executor Command, command string, args ...string) error {
	return executor.Execute(ctx, command, args...)
}
