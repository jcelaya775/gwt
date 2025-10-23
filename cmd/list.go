package cmd

import (
	"github.com/jcelaya775/gwt/git"
	"github.com/spf13/cobra"
)

func init() {
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
			println(wt)
		}

		return nil
	},
}
