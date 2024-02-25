package runner

import (
	"fmt"
	"path/filepath"

	"github.com/c-128/auroraci/internal/pipelines/artifacts"
	"github.com/c-128/auroraci/internal/pipelines/logger"
	"github.com/go-git/go-billy/v5"
)

func uploadArtifacts(
	uploader artifacts.Uploader,
	logger logger.Logger,
	worktree billy.Filesystem,
	artifacts []string,
) error {

	err := artifactsDir(
		uploader,
		logger,
		worktree,
		artifacts,
		worktree.Root(),
		worktree.Root(),
	)
	if err != nil {
		return err
	}

	return nil
}

func artifactsDir(
	uploader artifacts.Uploader,
	logger logger.Logger,
	worktree billy.Filesystem,
	artifacts []string,
	dir string,
	root string,
) error {
	files, err := worktree.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		fullPath := worktree.Join(dir, file.Name())

		if file.IsDir() {
			err := artifactsDir(
				uploader,
				logger,
				worktree,
				artifacts,
				fullPath,
				root,
			)
			if err != nil {
				return err
			}

			continue
		}

		relPath, err := filepath.Rel(root, fullPath)
		if err != nil {
			return fmt.Errorf("failed to rel path: %w", err)
		}

		for _, artifact := range artifacts {
			matched, err := filepath.Match(artifact, relPath)
			if err != nil {
				continue
			}

			if !matched {
				continue
			}

			logger.Printf("Uploading artifact \"%s\"", file.Name())
			reader, err := worktree.Open(fullPath)
			if err != nil {
				logger.Printf("Failed to upload artifact \"%s\": %s", file.Name(), err)
				continue
			}
			defer reader.Close()

			err = uploader.Upload(file.Name(), reader)
			if err != nil {
				logger.Printf("Failed to upload artifact \"%s\": %s", file.Name(), err)
				continue
			}
		}
	}

	return nil
}
