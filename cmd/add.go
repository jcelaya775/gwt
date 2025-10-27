package cmd

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/jcelaya775/gwt/internal/config"
	"github.com/jcelaya775/gwt/internal/connector"
	"github.com/jcelaya775/gwt/internal/git"
	"github.com/jcelaya775/gwt/internal/selecter"
	"github.com/jcelaya775/gwt/internal/shell"
	"github.com/jcelaya775/gwt/internal/zoxide"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func Add(c *config.Config, g *git.Git, s *selecter.Select, z *zoxide.Zoxide, conn *connector.Connector, sh shell.Shell) *cobra.Command {
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

			worktreePath, err := g.AddWorktree(c, branch, commitish, noPull, forceAdd)
			if err != nil {
				return err
			}
			boldStyle := lipgloss.NewStyle().Bold(true)
			fmt.Printf("Worktree for branch '%s' added successfully.\n\n", boldStyle.Render(branch))

			err = z.AddPath(worktreePath)
			if err != nil {
				return err
			}
			if seshConnect {
				err = conn.SeshConnect(worktreePath)
				if err != nil {
					return err
				}
			}

			if webStormConnect {
				err = conn.WebstormConnect(worktreePath)
				if err != nil {
					return err
				}
			}
			if goLandConnect {
				err = conn.GoLandConnect(worktreePath)
				if err != nil {
					return err
				}
			}
			if pyCharmConnect {
				err = conn.PyCharmConnect(worktreePath)
				if err != nil {
					return err
				}
			}
			if intelliJConnect {
				err = conn.IntelliJConnect(worktreePath)
				if err != nil {
					return err
				}
			}
			if cLionConnect {
				err = conn.CLionConnect(worktreePath)
				if err != nil {
					return err
				}
			}
			if riderConnect {
				err = conn.RiderConnect(worktreePath)
				if err != nil {
					return err
				}
			}
			if dataGripConnect {
				err = conn.DataGripConnect(worktreePath)
				if err != nil {
					return err
				}
			}

			var styledCommand string
			greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
			for i, command := range c.InitCommands {
				var execCmd *exec.Cmd
				if seshConnect {
					session := filepath.Base(worktreePath)
					execCmd = exec.Command("tmux", "send-keys", "-t", session, command, "C-m")
					execCmd.Stdout = os.Stdout
					execCmd.Stderr = os.Stderr
					execCmd.Stdin = os.Stdin

					commandText := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Render(command)
					styledCommand = greenStyle.Render(fmt.Sprintf("tmux send-keys -t %s ", session)) +
						commandText + greenStyle.Render(" C-m")
				} else {
					execCmd = exec.Command("sh", "-c", command)
					execCmd.Dir = worktreePath
					execCmd.Stdout = os.Stdout
					execCmd.Stderr = os.Stderr
					execCmd.Stdin = os.Stdin

					styledCommand = greenStyle.Render(command)
				}

				fmt.Println(boldStyle.Render(fmt.Sprintf("️➡️ Running init command %d of %s: %s...",
					i+1, strconv.Itoa(len(c.InitCommands)), styledCommand)))
				if err := execCmd.Run(); err != nil {
					return fmt.Errorf("error running init command '%s': %w", command, err)
				}
			}

			return nil
		},
	}

	addCmd.Flags().BoolVar(&noPull, "no-pull", false, "Do not pull the base branch before creating the worktree")
	addCmd.Flags().BoolVar(&noSync, "no-sync", false, "Do not fetch remote branches before creating the worktree")
	addCmd.Flags().BoolVarP(&forceAdd, "force", "f", false, "Checkout branch even if already checked out in another worktree")
	addCmd.Flags().BoolVar(&seshConnect, "sesh", false, "Connect to the worktree with seshConnect")
	addCmd.Flags().BoolVar(&webStormConnect, "webstorm", false, "Open the new worktree in WebStorm")
	addCmd.Flags().BoolVar(&intelliJConnect, "idea", false, "Open the new worktree in IntelliJ IDEA")
	addCmd.Flags().BoolVar(&pyCharmConnect, "pycharm", false, "Open the new worktree in PyCharm")
	addCmd.Flags().BoolVar(&cLionConnect, "clion", false, "Open the new worktree in CLion")
	addCmd.Flags().BoolVar(&riderConnect, "rider", false, "Open the new worktree in Rider")
	addCmd.Flags().BoolVar(&goLandConnect, "goland", false, "Open the new worktree in GoLand")
	addCmd.Flags().BoolVar(&dataGripConnect, "datagrip", false, "Open the new worktree in DataGrip")

	return addCmd
}
