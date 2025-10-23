package cmd

import (
	"fmt"
	"github.com/jcelaya775/gwt/git"
	"github.com/spf13/cobra"
)

var forceRemove bool
var keepBranch bool

func init() {
	removeCmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Force removal of the worktree even if there are uncommitted changes")
	removeCmd.Flags().BoolVarP(&keepBranch, "keep-branch", "k", false, "Also delete the branch associated with the worktree")
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:     "remove <worktree>",
	Short:   "Remove a git worktree",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		g, _ := git.New()
		worktrees, err := g.ListWorktrees()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return worktrees, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		worktree := args[0]

		g, err := git.New()
		if err != nil {
			return err
		}
		if err := g.RemoveWorktree(worktree, forceRemove, keepBranch); err != nil {
			return err
		}

		fmt.Printf("Worktree '%s' removed successfully.\n", worktree)
		return nil
	},
}
