package git

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/huh/spinner"
	"os/exec"
	"path/filepath"
	"strings"
)

type Git struct {
	worktreeRoot string
}

func New() (*Git, error) {
	worktreeRoot, err := getWorktreeRoot()
	if err != nil {
		return nil, err
	}
	return &Git{worktreeRoot: worktreeRoot}, nil
}

func (g *Git) GetWorktreeRoot() string {
	return g.worktreeRoot
}

func (g *Git) GetRepoName() string {
	return filepath.Base(g.worktreeRoot)
}

func (g *Git) AddWorktree(opts AddWorktreeOptions) error {
	cmdArgs := []string{"-C", g.worktreeRoot, "worktree", "add"}
	if opts.Branch == "" {
		// TODO: If worktreeRoot/commit is not provided, fzf/huh select from list of local and remote branches (avoid duplicates)
		fmt.Println("Not yet implemented: branch selection UI")
		return nil
	} else {
		existsLocally, err := g.BranchExistsLocally(opts.Branch)
		if err != nil {
			return err
		}
		existsRemotely, err := g.BranchExistsRemotely(opts.Branch)
		if err != nil {
			return err
		}
		if opts.Commitish != "" {
			cmdArgs = append(cmdArgs, "-b", opts.Branch, opts.Branch, opts.Commitish)
		} else {
			// Assume opts.Branch is opts.Branch name, and does not include a relative or absolute opts.Branch (TODO: validate)
			if existsLocally {
				cmdArgs = append(cmdArgs, opts.Branch, "--checkout", opts.Branch)
			} else if existsRemotely {
				cmdArgs = append(cmdArgs, "-b", opts.Branch, opts.Branch, "--checkout", fmt.Sprintf("origin/%s", opts.Branch))
			} else {
				cmdArgs = append(cmdArgs, "-b", opts.Branch, opts.Branch, "main") // TODO: detect default opts.Branch from config ⭐
			}
		}
	}

	output, err := exec.Command("git", cmdArgs...).CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}
	return nil
}

func (g *Git) ListWorktrees() ([]string, error) {
	output, err := exec.Command("git", "worktree", "list").Output()
	if err != nil {
		return nil, errors.New(string(output))
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	worktrees := make([]string, 0, len(lines))
	for _, line := range lines[1:] {
		worktreeAbsPath := strings.Fields(line)[0]
		worktreeRelPath := strings.TrimPrefix(worktreeAbsPath, g.worktreeRoot)
		worktrees = append(worktrees, worktreeRelPath)
	}
	return worktrees, nil
}

func (*Git) Fetch() error {
	var err error
	_ = spinner.New().
		Title("Syncing with remote branches...").
		Action(func() {
			cmd := exec.Command("git", "fetch")
			output, innerErr := cmd.CombinedOutput()
			if innerErr != nil {
				err = errors.New(string(output))
			}
		}).
		Run()
	return err
}

func (*Git) BranchExistsLocally(branch string) (bool, error) {
	output, err := exec.Command("git", "branch", "--list", branch).CombinedOutput()
	if err != nil {
		return false, errors.New(string(output))
	}

	if strings.TrimSpace(string(output)) == branch {
		return true, nil
	}
	return false, nil
}

func (*Git) BranchExistsRemotely(branch string) (bool, error) {
	remoteBranch := "origin/" + branch
	output, err := exec.Command("git", "branch", "-r", "--list", remoteBranch).CombinedOutput()
	if err != nil {
		return false, errors.New(string(output))
	}

	if strings.TrimSpace(string(output)) == remoteBranch {
		return true, nil
	}
	return false, nil
}

func getWorktreeRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-common-dir")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.New("not in a git repository")
	}

	gitDir := strings.TrimSpace(string(output))
	if strings.HasSuffix(gitDir, ".git") {
		return fmt.Sprintf("%s/", filepath.Dir(gitDir)), nil
	} else {
		return "", errors.New("could not find git repository root containing .git. Please use gwt clone to clone the repository")
	}
}
