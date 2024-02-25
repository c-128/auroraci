package runner

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/c-128/auroraci/internal/tarball"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
)

func copyWorktreeIntoContainer(
	ctx context.Context,
	docker *client.Client,
	containerID string,
	worktree billy.Filesystem,
	target string,
) error {
	reader, writer := io.Pipe()
	errChannel := make(chan error)

	go func() {
		defer writer.Close()

		err := tarball.FSToTarball(
			worktree,
			worktree.Root(),
			writer,
		)
		if err != nil {
			errChannel <- fmt.Errorf("failed to create tarball: %w", err)
			return
		}

		errChannel <- nil
	}()

	go func() {
		defer reader.Close()

		err := docker.CopyToContainer(
			ctx,
			containerID,
			target,
			reader,
			types.CopyToContainerOptions{},
		)
		if err != nil {
			errChannel <- fmt.Errorf("failed to copy to container: %w", err)
			return
		}

		errChannel <- nil
	}()

	err := errors.Join(<-errChannel, <-errChannel)
	if err != nil {
		return err
	}

	return nil
}

func copyContainerIntoWorktree(
	ctx context.Context,
	docker *client.Client,
	containerID string,
	target string,
) (billy.Filesystem, error) {
	reader, _, err := docker.CopyFromContainer(
		ctx,
		containerID,
		target,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to copy from container: %w", err)
	}
	defer reader.Close()

	worktree := memfs.New()
	err = tarball.TarballToFS(
		reader,
		worktree.Root(),
		worktree,
		target,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to extract tarball: %w", err)
	}

	return worktree, nil
}
