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

func Add(c *config.Config, g *git.Git, s *selecter.Select) *cobra.Command {
	var noPull bool
	var noSync bool
	var forceAdd bool

	addCmd := &cobra.Command{
		Use:     "add <branch> [commit-ish]",
		Short:   "Add a new worktree",
		Aliases: []string{"a"},
		Args:    cobra.MaximumNArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			if len(args) > 1 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			branches, err := g.ListBranches(false, true)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return branches, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var branch, commitish string

			if !noSync {
				if err := g.Fetch(); err != nil {
					return err
				}
			}

			if len(args) == 0 {
				branchesToSelectFrom, err := g.ListBranches(false, true)
				if err != nil {
					return err
				}

				if len(branchesToSelectFrom) == 0 {
					return errors.New("no branches available to select from")
				}

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

			err = g.AddWorktree(c, branch, commitish, noPull, forceAdd)
			if err != nil {
				return err
			}
			fmt.Printf("Worktree for branch '%s' added successfully.\n", branch)
			return nil
		},
	}

	addCmd.Flags().BoolVar(&noPull, "no-pull", false, "Do not pull the base branch before creating the worktree")
	addCmd.Flags().BoolVar(&noSync, "no-sync", false, "Do not fetch remote branches before creating the worktree")
	addCmd.Flags().BoolVarP(&forceAdd, "force", "f", false, "Checkout branch even if already checked out in another worktree")

	return addCmd
}
