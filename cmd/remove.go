package cmd

import (
	"fmt"
	"github.com/jcelaya775/gwt/internal/git"
	"github.com/jcelaya775/gwt/internal/selecter"
	"github.com/spf13/cobra"
)

var forceRemove bool
var keepBranch bool

func Remove(g *git.Git, s *selecter.Select) *cobra.Command {
	removeCmd := &cobra.Command{
		Use:     "remove [worktree...]",
		Short:   "Remove a git worktree",
		Aliases: []string{"rm"},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			worktrees, err := g.ListWorktrees()
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return worktrees, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			worktrees := args

			if len(worktrees) == 0 {
				availableWorktrees, err := g.ListWorktrees()
				worktrees, err = s.MultiSelect("Select worktrees to remove:", availableWorktrees)
				if err != nil {
					return err
				}
			}

			for _, worktree := range worktrees {
				if err := g.RemoveWorktree(worktree, forceRemove, keepBranch); err != nil {
					return err
				}
				fmt.Printf("Worktree '%s' removed successfully.\n", worktree)
			}

			return nil
		},
	}

	removeCmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Force removal of the worktree even if there are uncommitted changes")
	removeCmd.Flags().BoolVarP(&keepBranch, "keep-branch", "k", false, "Also delete the branch associated with the worktree")

	return removeCmd
}
