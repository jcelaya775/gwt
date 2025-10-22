package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Git interface {
	BranchExistsLocally(branch string) (bool, error)
	BranchExistsRemotely(branch string) (bool, error)
}

type RealGit struct{}

func New() Git {
	return &RealGit{}
}

func (*RealGit) BranchExistsLocally(branch string) (bool, error) {
	output, err := exec.Command("git", "branch", "--list", branch).CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("%s", string(output))
	}

	if strings.TrimSpace(string(output)) == branch {
		return true, nil
	}
	return false, nil
}

func (*RealGit) BranchExistsRemotely(branch string) (bool, error) {
	remoteBranch := "origin/" + branch
	output, err := exec.Command("git", "branch", "-r", "--list", remoteBranch).CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("%s", string(output))
	}

	if strings.TrimSpace(string(output)) == remoteBranch {
		return true, nil
	}
	return false, nil
}
