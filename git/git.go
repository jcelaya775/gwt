package git

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/huh/spinner"
	"github.com/jcelaya775/gwt/internal/config"
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

func (g *Git) AddWorktree(config *config.Config, branch string, baseBranch string, commitish string, pull bool, force bool) error {
	cmdArgs := []string{"-C", g.worktreeRoot, "worktree", "add"}

	if branch == "" {
		// TODO: If worktreeRoot/commit is not provided, fzf/huh select from list of local and remote branches (avoid duplicates)
		fmt.Println("Not yet implemented: branch selection UI")
		return nil
	} else {
		existsLocally, err := g.BranchExistsLocally(branch)
		if err != nil {
			return err
		}
		existsRemotely, err := g.BranchExistsRemotely(branch)
		if err != nil {
			return err
		}

		if commitish != "" {
			cmdArgs = append(cmdArgs, "-b", branch, branch, commitish)
		} else {
			// TODO: Add base branch flag (use config.Defaults.BaseBranch as default)
			if existsLocally {
				cmdArgs = append(cmdArgs, branch, "--checkout", branch)
			} else if existsRemotely {
				cmdArgs = append(cmdArgs, "-b", branch, branch, "--checkout", fmt.Sprintf("origin/%s", branch))
			} else {
				// TODO: What if user wants to branch off of current branch? Add flag to specify base branch
				cmdArgs = append(cmdArgs, "-b", branch, branch, config.Defaults.BaseBranch)
			}
		}
	}

	output, err := exec.Command("git", cmdArgs...).CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}
	return nil
}

func (g *Git) RemoveWorktree(worktree string, force, keepBranch bool) error {
	branch, err := g.GetWorktreeBranch(worktree)

	cmdArgs := []string{"-C", g.worktreeRoot, "worktree", "remove"}
	if force {
		cmdArgs = append(cmdArgs, "--force")
	}
	cmdArgs = append(cmdArgs, worktree)

	output, err := exec.Command("git", cmdArgs...).CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}

	if !keepBranch {
		if err != nil {
			return err
		}
		branchOutput, branchErr := exec.Command("git", "-C", g.worktreeRoot, "branch", "-D", branch).CombinedOutput()
		if branchErr != nil {
			return errors.New(string(branchOutput))
		}
	}

	return nil
}

func (g *Git) GetWorktreeBranch(worktree string) (string, error) {
	cmdArgs := []string{"-C", g.worktreeRoot, "worktree", "list"}
	output, err := exec.Command("git", cmdArgs...).Output()
	if err != nil {
		return "", errors.New(string(output))
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		worktreeAbsPath := fields[0]
		worktreeRelPath := strings.TrimPrefix(worktreeAbsPath, g.worktreeRoot)
		if worktreeRelPath == worktree {
			if len(fields) >= 3 {
				unparsedBranch := fields[2]
				return strings.TrimSuffix(strings.TrimPrefix(unparsedBranch, "["), "]"), nil
			}
			return "", errors.New("could not determine branch for worktree")
		}
	}
	return "", errors.New("worktree not found")
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
		gitDirAbsolute, err := filepath.Abs(gitDir)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s/", filepath.Dir(gitDirAbsolute)), nil
	} else {
		return "", errors.New("could not find git repository root containing .git. Please use gwt clone to clone the repository")
	}
}
