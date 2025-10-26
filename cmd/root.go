/*
Copyright © 2025 Jorge Celaya jcelaya775@gmail.com
*/
package cmd

import (
	"github.com/jcelaya775/gwt/internal/config"
	"github.com/jcelaya775/gwt/internal/git"
	"github.com/jcelaya775/gwt/internal/selecter"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gwt",
	Short: "A git worktree wrapper that makes life easier\n\n",
	/// TODO: Show TUI when no subcommand is provided
}

func Execute() {
	log.SetFlags(0)

	g, err := git.New()
	if err != nil {
		log.Fatalf("failed to initialize git: %v", err)
	}
	c, err := config.LoadConfig(g.GetWorktreeRoot())
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	s := selecter.New()

	rootCmd.AddCommand(Add(c, g, s))
	rootCmd.AddCommand(Clone(g))
	rootCmd.AddCommand(List(g))
	rootCmd.AddCommand(Remove(g, s))

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
