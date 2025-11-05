/*
Copyright © 2025 Jorge Celaya jcelaya775@gmail.com
*/
package cmd

import (
	_connector "github.com/jcelaya775/gwt/internal/connector"
	_git "github.com/jcelaya775/gwt/internal/git"
	_home "github.com/jcelaya775/gwt/internal/home"
	_selecter "github.com/jcelaya775/gwt/internal/selecter"
	_sesh "github.com/jcelaya775/gwt/internal/sesh"
	_shell "github.com/jcelaya775/gwt/internal/shell"
	_zoxide "github.com/jcelaya775/gwt/internal/zoxide"
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

	git, err := _git.New()
	if err != nil {
		log.Fatalf("failed to initialize git: %v", err)
	}
	selecter := _selecter.New()
	home := _home.NewHome()
	shell := _shell.NewShell(home)
	zoxide := _zoxide.New(shell)
	connector := _connector.New(shell)
	sesh := _sesh.New(shell)

	rootCmd.AddCommand(Add(git, selecter, zoxide, connector))
	rootCmd.AddCommand(Clone(git))
	rootCmd.AddCommand(List(git))
	rootCmd.AddCommand(Remove(git, sesh, selecter))
	rootCmd.AddCommand(Init(git))

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
