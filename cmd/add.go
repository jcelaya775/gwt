package cmd

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	_config "github.com/jcelaya775/gwt/internal/config"
	"github.com/jcelaya775/gwt/internal/connector"
	"github.com/jcelaya775/gwt/internal/git"
	"github.com/jcelaya775/gwt/internal/selecter"
	"github.com/jcelaya775/gwt/internal/utils"
	"github.com/jcelaya775/gwt/internal/zoxide"
	"github.com/spf13/cobra"
	"strings"
)

func Add(git *git.Git, selecter *selecter.Select, zoxide *zoxide.Zoxide, connector *connector.Connector) *cobra.Command {
	var noPull bool
	var noSync bool
	var forceAdd bool
	var seshConnect bool
	var webStormConnect bool
	var goLandConnect bool
	var pyCharmConnect bool
	var intelliJConnect bool
	var cLionConnect bool
	var riderConnect bool
	var dataGripConnect bool

	addCmd := &cobra.Command{
		Use:     "add <branch> [commit-ish]",
		Short:   "Add a new worktree",
		Aliases: []string{"a"},
		Args:    cobra.MaximumNArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			if len(args) > 1 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			err := git.SetWorktreeRoot()
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			branches, err := git.ListBranches(false, true)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return branches, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var branch, commitish string

			if git.GetWorktreeRoot() == "" {
				if err := git.SetWorktreeRoot(); err != nil {
					return err
				}
			}

			config, err := _config.LoadConfig(git.GetWorktreeRoot())
			if err != nil {
				return err
			}

			if !noSync {
				if err := git.Fetch(); err != nil {
					return err
				}
			}

			if len(args) == 0 {
				branchesToSelectFrom, err := git.ListBranches(false, true)
				if err != nil {
					return err
				}

				if len(branchesToSelectFrom) == 0 {
					return errors.New("no branches available to select from")
				}

				branch, err = selecter.Select("Select a branch to create a worktree:", branchesToSelectFrom)
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

			worktreeAlreadyExists, err := git.WorktreeExists(strings.TrimPrefix(branch, "origin/"))
			if err != nil {
				return err
			}
			if worktreeAlreadyExists {
				return fmt.Errorf("worktree for branch '%s' already exists", branch)
			}

			worktreePath, err := git.AddWorktree(config, branch, commitish, noPull, forceAdd)
			if err != nil {
				return err
			}
			boldStyle := lipgloss.NewStyle().Bold(true)
			fmt.Printf("Worktree for branch %s added successfully.\n\n", boldStyle.Render(branch))

			if err = zoxide.AddPath(worktreePath); err != nil {
				return err
			}
			if seshConnect {
				if err = connector.SeshConnect(worktreePath); err != nil {
					return err
				}
			}

			if webStormConnect {
				if err = connector.WebstormConnect(worktreePath); err != nil {
					return err
				}
			}
			if goLandConnect {
				if err = connector.GoLandConnect(worktreePath); err != nil {
					return err
				}
			}
			if pyCharmConnect {
				if err = connector.PyCharmConnect(worktreePath); err != nil {
					return err
				}
			}
			if intelliJConnect {
				if err = connector.IntelliJConnect(worktreePath); err != nil {
					return err
				}
			}
			if cLionConnect {
				if err = connector.CLionConnect(worktreePath); err != nil {
					return err
				}
			}
			if riderConnect {
				if err = connector.RiderConnect(worktreePath); err != nil {
					return err
				}
			}
			if dataGripConnect {
				if err = connector.DataGripConnect(worktreePath); err != nil {
					return err
				}
			}

			if err = utils.RunCommands(config.InitCommands, worktreePath, seshConnect, ""); err != nil {
				return err
			}

			return nil
		},
	}

	addCmd.Flags().BoolVar(&noPull, "no-pull", false, "Do not pull the base branch before creating the worktree")
	addCmd.Flags().BoolVar(&noSync, "no-sync", false, "Do not fetch remote branches before creating the worktree")
	addCmd.Flags().BoolVarP(&forceAdd, "force", "f", false, "Checkout branch even if already checked out in another worktree")
	addCmd.Flags().BoolVar(&seshConnect, "sesh", false, "Connect to the worktree with sesh")
	addCmd.Flags().BoolVar(&webStormConnect, "webstorm", false, "Open the new worktree in WebStorm")
	addCmd.Flags().BoolVar(&intelliJConnect, "idea", false, "Open the new worktree in IntelliJ IDEA")
	addCmd.Flags().BoolVar(&pyCharmConnect, "pycharm", false, "Open the new worktree in PyCharm")
	addCmd.Flags().BoolVar(&cLionConnect, "clion", false, "Open the new worktree in CLion")
	addCmd.Flags().BoolVar(&riderConnect, "rider", false, "Open the new worktree in Rider")
	addCmd.Flags().BoolVar(&goLandConnect, "goland", false, "Open the new worktree in GoLand")
	addCmd.Flags().BoolVar(&dataGripConnect, "datagrip", false, "Open the new worktree in DataGrip")

	return addCmd
}
