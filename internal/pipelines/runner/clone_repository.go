package runner

import (
	"fmt"

	"github.com/c-128/auroraci/internal/pipelines"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

func cloneRepository(repo *pipelines.Repository) (billy.Filesystem, error) {
	worktree := memfs.New()

	_, err := git.Clone(
		memory.NewStorage(),
		worktree,
		&git.CloneOptions{
			URL:           repo.Origin,
			ReferenceName: plumbing.NewBranchReferenceName(repo.Branch),
			SingleBranch:  true,
			Depth:         1,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	return worktree, nil
}
