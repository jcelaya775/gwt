package cmd

import (
	"fmt"
	"github.com/charmbracelet/huh/spinner"
	"github.com/jcelaya775/gwt/git"
	"github.com/spf13/cobra"
	"os/exec"
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
	Use:     "add <path> [commit-ish]",
	Short:   "Add a new worktree",
	Args:    cobra.MaximumNArgs(2),
	Aliases: []string{"h"},
	RunE: func(cmd *cobra.Command, args []string) error {
		git := git.New()

		if !noSync {
			var err error
			_ = spinner.New().
				Title("Syncing with remote branches...").
				Action(func() {
					var output []byte
					output, err = exec.Command("git", "fetch").CombinedOutput()
					if err != nil {
						err = fmt.Errorf("%s", string(output))
					}
				}).
				Run()
			if err != nil {
				return err
			}
		}

		// TODO: Find root repo path from current directory
		cmdArgs := []string{"worktree", "add"}
		if len(args) == 0 {
			// TODO: If path/commit is not provided, fzf select from list of local and remote branches (avoid duplicates)
		} else {
			var path, commitish string
			path = args[0]
			if len(args) == 2 {
				commitish = args[1]
			}

			if commitish == "" {
				// Assume path is branch name, and does not include a relative or absolute path (TODO: validate)
				existsLocally, err := git.BranchExistsLocally(path)
				if err != nil {
					return err
				}
				existsRemotely, err := git.BranchExistsRemotely(path)
				if err != nil {
					return err
				}
				if existsRemotely {
					cmdArgs = append(cmdArgs, "-b", path, path, "--checkout", fmt.Sprintf("origin/%s", path))
				} else if existsLocally {
					cmdArgs = append(cmdArgs, path)
				} else {
					cmdArgs = append(cmdArgs, "-b", path, path, baseBranch)
				}
			} else {
				// If commit-ish is provided, always create new branch with worktree name
				cmdArgs = append(cmdArgs, "-b", path, commitish)
			}
		}

		output, err := exec.Command("git", cmdArgs...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v", string(output))
		}
		return nil
	},
}
