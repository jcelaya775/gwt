package cmd

import (
	"fmt"
	"github.com/jcelaya775/gwt/git"
	"github.com/spf13/cobra"
	"path/filepath"
)

var absolutePath bool

func init() {
	listCmd.Flags().BoolVarP(&absolutePath, "absolute", "a", false, "Show absolute paths")
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all worktrees",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		g, err := git.New()
		if err != nil {
			return err
		}

		worktrees, err := g.ListWorktrees()
		if err != nil {
			return err
		}

		for _, wt := range worktrees {
			if absolutePath {
				fmt.Println(filepath.Join(g.GetWorktreeRoot(), wt))
			} else {
				fmt.Println(wt)
			}
		}

		return nil
	},
}
