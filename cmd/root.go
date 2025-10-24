/*
Copyright © 2025 Jorge Celaya jcelaya775@gmail.com
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gwt",
	Short: "A git worktree wrapper that makes life easier\n\n",
	/// TODO: Show TUI when no subcommand is provided
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// TODO: Add config path flag
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gwt.yaml)")
}
