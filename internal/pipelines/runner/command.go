package runner

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/c-128/auroraci/internal/pipelines"
	"github.com/c-128/auroraci/internal/pipelines/logger"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func executeCommand(
	ctx context.Context,
	logger logger.Logger,
	docker *client.Client,
	containerID string,
	command pipelines.BuildCommand,
) (int, error) {
	var cmd []string
	switch {
	case command.Run != "":
		cmd = strings.Split(command.Run, " ")
	case command.RunBash != "":
		cmd = []string{"/bin/bash", "-c", command.RunBash}
	}

	logger.Printf("Executing command \"%s\"", strings.Join(cmd, " "))

	exec, err := docker.ContainerExecCreate(
		ctx,
		containerID,
		types.ExecConfig{
			Cmd:          cmd,
			AttachStdout: true,
			AttachStderr: true,
		},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create container exec: %w", err)
	}

	attach, err := docker.ContainerExecAttach(
		ctx,
		exec.ID,
		types.ExecStartCheck{},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to attach to container exec: %w", err)
	}
	defer attach.Close()

	logger.Printf("------------------------------------------------------------------------------------------------")
	defer logger.Printf("------------------------------------------------------------------------------------------------")

	for {
		char, _, err := attach.Reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("failed to read from container exec: %w", err)
		}

		logger.PrintCommand(char)
	}

	inspected, err := docker.ContainerExecInspect(ctx, exec.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to inspect container exec: %w", err)
	}

	return inspected.ExitCode, nil
}
