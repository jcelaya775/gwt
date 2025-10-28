package cmd

import (
	"fmt"
	"github.com/jcelaya775/gwt/internal/git"
	"github.com/spf13/cobra"
	"path/filepath"
)

func List(g *git.Git) *cobra.Command {
	var absolutePath bool

	var listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List all worktrees",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := g.SetWorktreeRoot()
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

	listCmd.Flags().BoolVarP(&absolutePath, "absolute", "a", false, "Show absolute paths")

	return listCmd
}
