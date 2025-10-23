package cmd

import (
	"fmt"
	"github.com/jcelaya775/gwt/git"
	"github.com/spf13/cobra"
)

// var newBranch string
// var checkout bool
var pull bool
var noSync bool
var force bool
var baseBranch = "main"

func init() {
	// TODO: Assume new branch if branch is not found locally or remotely
	//cloneCmd.Flags().StringVar(&newBranch, "b", "", "Create a new branch")
	//cloneCmd.Flags().BoolVar(&checkout, "checkout", true, "Checkout branch in the new worktree. Sets remote tracking if <branch> exists on remote")
	addCmd.Flags().BoolVarP(&pull, "pull", "p", false, "Pull the base branch before creating the worktree")
	addCmd.Flags().BoolVar(&noSync, "no-sync", false, "Do not fetch remote branches before creating the worktree")
	addCmd.Flags().BoolVarP(&force, "force", "f", false, "Checkout branch even if already checked out in another worktree")
	rootCmd.AddCommand(addCmd)
}

// TODO: Add config support with default base branch and pull options
var addCmd = &cobra.Command{
	Use:   "add <branch> [commit-ish]",
	Short: "Add a new worktree",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		g, err := git.New()
		if err != nil {
			return err
		}

		if !noSync {
			if err := g.Fetch(); err != nil {
				return err
			}
		}

		if len(args) == 0 {
			// TODO: If path/commit is not provided, fzf/huh select from list of local and remote branches (avoid duplicates)
			fmt.Println("Not yet implemented: branch selection UI")
		} else {
			var branch, commitish string
			branch = args[0]
			if len(args) == 2 {
				commitish = args[1]
			}

			err = g.AddWorktree(git.AddWorktreeOptions{
				Branch:    branch,
				Commitish: commitish,
				Pull:      pull,
				Force:     force,
			})
			if err != nil {
				return err
			}
		}

		// TODO: Display success  message with created worktree path and/or cd/sesh into it & open IDE (option)
		return nil
	},
}
