package cmd

import (
	"github.com/jcelaya775/gwt/internal/git"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cloneCmd)
}

var cloneCmd = &cobra.Command{
	Use:   "clone <repo> [dir]",
	Short: "Clone a git repository in a worktree setup",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var repoURL, dir string
		repoURL = args[0]
		if len(args) == 2 {
			dir = args[1]
		}

		var g git.Git
		err := g.CloneRepo(repoURL, dir)
		if err != nil {
			return err
		}

		return nil
	},
}
