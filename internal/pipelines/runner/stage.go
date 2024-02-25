package runner

import (
	"context"
	"fmt"

	"github.com/c-128/auroraci/internal/pipelines"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func createContainer(
	ctx context.Context,
	docker *client.Client,
	stage pipelines.BuildStage,
) (string, error) {
	_, _, err := docker.ImageInspectWithRaw(
		ctx,
		stage.Image,
	)
	if client.IsErrNotFound(err) {
		reader, err := docker.ImagePull(
			ctx,
			stage.Image,
			types.ImagePullOptions{},
		)
		if err != nil {
			return "", fmt.Errorf("failed to pull image: %w", err)
		}
		defer reader.Close()

		buf := make([]byte, 256)
		for {
			_, err := reader.Read(buf)
			if err != nil {
				break
			}
		}
	}

	createRes, err := docker.ContainerCreate(
		ctx,
		&container.Config{
			Image:      stage.Image,
			WorkingDir: stage.Workdir,
			Tty:        true,
		},
		&container.HostConfig{
			AutoRemove: true,
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	err = docker.ContainerStart(
		ctx,
		createRes.ID,
		container.StartOptions{},
	)
	if err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	return createRes.ID, nil
}

func removeContainer(
	ctx context.Context,
	docker *client.Client,
	id string,
) error {
	err := docker.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	err = docker.ContainerRemove(ctx, id, container.RemoveOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	return nil
}
