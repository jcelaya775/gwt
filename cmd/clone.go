package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"os/exec"
	"strings"
)

func init() {
	rootCmd.AddCommand(cloneCmd)
}

var cloneCmd = &cobra.Command{
	Use:   "clone <repo> [dir]",
	Short: "Clone a git repository in a worktree setup",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var repo, dir string
		repo = args[0]
		if len(args) == 2 {
			dir = args[1]
		}

		var cloneArgs []string
		if dir != "" {
			cloneArgs = []string{"clone", "--no-checkout", repo, dir}
		} else {
			cloneArgs = []string{"clone", "--no-checkout", repo}
		}
		output, err := exec.Command("git", cloneArgs...).CombinedOutput()
		if err != nil {
			return errors.New(string(output))
		}

		// TODO: Checkout a dummy branch and create worktree for main branch relative to worktree dir (config)
		// TODO: Create pseudo-random branch name to avoid conflicts, random 6 digit number
		var repoDir string
		if dir != "" {
			repoDir = dir
		} else {
			parts := strings.Split(repo, "/")
			repoName := parts[len(parts)-1]
			repoDir = strings.TrimSuffix(repoName, ".git")
		}

		output, err = exec.Command("git", "-C", repoDir, "branch", "--show-current").CombinedOutput()
		if err != nil {
			return errors.New(string(output))
		}
		originalBranch := strings.TrimSpace(string(output))

		dummyBranch := "dummy"
		output, err = exec.Command("git", "-C", repoDir, "checkout", "-b", dummyBranch).CombinedOutput()
		if err != nil {
			return errors.New(string(output))
		}

		var worktreeArgs []string
		if dir != "" {
			worktreeArgs = []string{"-C", dir, "worktree", "add", originalBranch, originalBranch}
		} else {
			worktreeArgs = []string{"worktree", "add", originalBranch, originalBranch}
		}
		output, err = exec.Command("git", worktreeArgs...).CombinedOutput()
		if err != nil {
			return errors.New(string(output))
		}
		return nil
	},
}
