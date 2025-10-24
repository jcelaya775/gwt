package cmd

import (
	"errors"
	"fmt"
	"github.com/jcelaya775/gwt/internal/config"
	"github.com/jcelaya775/gwt/internal/git"
	"github.com/jcelaya775/gwt/internal/selecter"
	"github.com/spf13/cobra"
	"strings"
)

var noPull bool
var noSync bool

func init() {
	addCmd.Flags().BoolVar(&noPull, "no-pull", false, "Do not pull the base branch before creating the worktree")
	addCmd.Flags().BoolVar(&noSync, "no-sync", false, "Do not fetch remote branches before creating the worktree")
	addCmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Checkout branch even if already checked out in another worktree")
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:     "add <branch> [commit-ish]",
	Short:   "Add a new worktree",
	Aliases: []string{"a"},
	Args:    cobra.MaximumNArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		g, err := git.New()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		branches, err := g.ListBranches(false, true)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return branches, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var branch, commitish string

		g, err := git.New()
		if err != nil {
			return err
		}

		c, err := config.LoadConfig(g.GetWorktreeRoot())
		if err != nil {
			return err
		}

		if len(args) == 0 {
			branchesToSelectFrom, err := g.ListBranches(false, true)
			if err != nil {
				return err
			}

			if len(branchesToSelectFrom) == 0 {
				return errors.New("no branches available to select from")
			}

			s := selecter.New()
			branch, err = s.Select("Select a branch to create a worktree:", branchesToSelectFrom)
			if err != nil {
				return err
			}
			if branch == "" {
				return nil
			}
		} else {
			branch = args[0]
		}
		if len(args) >= 2 {
			commitish = args[1]
		}

		worktreeAlreadyExists, err := g.WorktreeExists(strings.TrimPrefix(branch, "origin/"))
		if err != nil {
			return err
		}
		if worktreeAlreadyExists {
			return fmt.Errorf("worktree for branch '%s' already exists", branch)
		}

		if !noSync {
			if err := g.Fetch(); err != nil {
				return err
			}
		}

		err = g.AddWorktree(c, branch, commitish, noPull, forceRemove)
		if err != nil {
			return err
		}
		fmt.Printf("Worktree for branch '%s' added successfully.\n", branch)
		return nil
	},
}
