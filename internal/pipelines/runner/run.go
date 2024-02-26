package runner

import (
	"context"

	"github.com/c-128/auroraci/internal/pipelines"
	"github.com/c-128/auroraci/internal/pipelines/artifacts"
	"github.com/c-128/auroraci/internal/pipelines/logger"
	"github.com/docker/docker/client"
)

func RunPipeline(
	ctx context.Context,
	logger logger.Logger,
	uploader artifacts.Uploader,
	docker *client.Client,
	pipeline *pipelines.Pipeline,
) error {
	defer logger.Close()
	defer uploader.Close()

	logger.Printf("Cloning repository \"%s\"", pipeline.Repository.Origin)
	worktree, err := cloneRepository(pipeline.Repository)
	if err != nil {
		logger.Printf("Failed to clone repository \"%s\": %s", pipeline.Repository.Origin, err)
		return err
	}

	for _, stage := range pipeline.Build.Stages {
		logger.Printf("Preparing container")
		containerID, err := createContainer(ctx, docker, stage)
		if err != nil {
			logger.Fatalf("Failed to prepare container: %s", err)
			return err
		}
		logger.Printf("Container created with ID \"%s\"", containerID)

		defer logger.Printf("Cleaning up container \"%s\"", containerID)
		defer removeContainer(ctx, docker, containerID)

		logger.Printf("Copying worktree into container")
		err = copyWorktreeIntoContainer(ctx, docker, containerID, worktree, stage.Workdir)
		if err != nil {
			logger.Fatalf("Failed to copy worktree into container: %s", err)
			return err
		}

		for _, command := range stage.Commands {
			exitCode, err := executeCommand(ctx, logger, docker, containerID, command)
			if err != nil {
				logger.Fatalf("Failed to execute command: %s", err)
				return err
			}

			if exitCode != command.ExitCode {
				logger.Fatalf("Command exited with code %d instead of %d", exitCode, command.ExitCode)
				return nil
			}
		}

		logger.Printf("Copying worktree from container")
		worktree, err = copyContainerIntoWorktree(ctx, docker, containerID, stage.Workdir)
		if err != nil {
			logger.Fatalf("Failed to copy worktree from container: %s", err)
			return err
		}
	}

	err = uploadArtifacts(uploader, logger, worktree, pipeline.Build.Artifacts)
	if err != nil {
		logger.Fatalf("Failed to upload artifacts: %s", err)
		return err
	}

	return nil
}
