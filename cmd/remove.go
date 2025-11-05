package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	_config "github.com/jcelaya775/gwt/internal/config"
	"github.com/jcelaya775/gwt/internal/git"
	"github.com/jcelaya775/gwt/internal/selecter"
	"github.com/jcelaya775/gwt/internal/sesh"
	"github.com/jcelaya775/gwt/internal/utils"
	"github.com/spf13/cobra"
	"path/filepath"
)

var forceRemove bool
var keepBranch bool

func Remove(git *git.Git, sesh *sesh.Sesh, selecter *selecter.Select) *cobra.Command {
	removeCmd := &cobra.Command{
		Use:     "remove [worktree...]",
		Short:   "Remove a git worktree",
		Aliases: []string{"rm"},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			err := git.SetWorktreeRoot()
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			worktrees, err := git.ListWorktrees()
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return worktrees, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			worktrees := args

			if git.GetWorktreeRoot() == "" {
				err := git.SetWorktreeRoot()
				if err != nil {
					return err
				}
			}

			config, err := _config.LoadConfig(git.GetWorktreeRoot())
			if err != nil {
				return err
			}

			if len(worktrees) == 0 {
				availableWorktrees, err := git.ListWorktrees()
				worktrees, err = selecter.MultiSelect("Select worktrees to remove:", availableWorktrees)
				if err != nil {
					return err
				}
			}

			for i, worktree := range worktrees {
				boldStyle := lipgloss.NewStyle().Bold(true)
				worktreePath := filepath.Join(git.GetWorktreeRoot(), worktree)
				if err := utils.RunCommands(config.DestroyCommands, worktreePath, false, worktree); err != nil {
					return err
				}
				if len(config.DestroyCommands) > 0 {
					fmt.Println()
				}

				if err := git.RemoveWorktree(worktree, forceRemove, keepBranch); err != nil {
					return err
				}
				fmt.Printf("Worktree %s removed successfully.\n", boldStyle.Render(worktree))

				if i < len(worktrees)-1 {
					fmt.Println()
				}
			}

			return nil
		},
	}

	removeCmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Force removal of the worktree even if there are uncommitted changes")
	removeCmd.Flags().BoolVarP(&keepBranch, "keep-branch", "k", false, "Also delete the branch associated with the worktree")

	return removeCmd
}
