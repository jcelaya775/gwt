package cmd

import (
	"github.com/spf13/cobra"
	"os/exec"
	"strings"
)

func init() {
	rootCmd.AddCommand(cloneCmd)
}

var cloneCmd = &cobra.Command{
	Use:     "clone <repo> [dir]",
	Short:   "Clone a git repository that uses worktrees",
	Args:    cobra.RangeArgs(1, 2),
	Aliases: []string{"h"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var repo, dir string
		repo = args[0]
		if len(args) == 2 {
			dir = args[1]
		}

		err := exec.Command("git", "clone", repo, dir, "--no-checkout").Run()
		if err != nil {
			return err
		}

		// TODO: Checkout a dummy branch and create worktree for main branch relative to worktree dir (config)
		// TODO: Create pseudo-random branch name to avoid conflicts, random 6 digit number
		output, err := exec.Command("git", "branch", "--show-current").Output()
		if err != nil {
			return err
		}
		originalBranch := strings.TrimSpace(string(output))

		var checkoutArgs []string
		dummyBranch := "dummy"
		if dir != "" {
			checkoutArgs = []string{"-C", dir, "checkout", "-b", dummyBranch}
		} else {
			checkoutArgs = []string{"checkout", "-b", dummyBranch}
		}
		err = exec.Command("git", checkoutArgs...).Run()
		if err != nil {
			return err
		}

		var worktreeArgs []string
		if dir != "" {
			worktreeArgs = []string{"-C", dir, "worktree", "add", originalBranch, originalBranch}
		} else {
			worktreeArgs = []string{"worktree", "add", originalBranch, originalBranch}
		}
		return exec.Command("git", worktreeArgs...).Run()
	},
}
