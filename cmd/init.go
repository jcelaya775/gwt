package cmd

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/jcelaya775/gwt/internal/git"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func Init(g *git.Git) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize gwt configuration in the current git repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := g.SetWorktreeRoot()
			if err != nil {
				return err
			}
			configPath := filepath.Join(g.GetWorktreeRoot(), ".gwt.yml")
			if _, err := os.Stat(configPath); err == nil {
				var confirm bool
				err := huh.NewConfirm().
					Title("Config file already exists. Do you want to overwrite it?").
					Affirmative("Yes").
					Negative("No").
					Value(&confirm).
					Run()
				if err != nil {
					return err
				}
				if !confirm {
					return nil
				}
			}

			configContent := `# gwt configuration
version: "1.0"

# Default settings for worktrees
defaults:
  # Default base branch for new worktrees
  base_branch: main

# Commands that run after creating a worktree
init_commands:
  - echo "Worktree initialized!"

# Commands that run when removing a worktree
remove_commands:
  - echo "Worktree removed!"
`
			if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
				return err
			}

			fmt.Println("Initialized gwt configuration at", configPath)
			return nil
		},
	}
}
