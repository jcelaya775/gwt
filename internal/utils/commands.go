package utils

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func RunCommands(commands []string, worktreePath string, sesh bool, worktree string) error {
	var styledCommandText string
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	orangeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500"))
	for i, command := range commands {
		var execCmd *exec.Cmd
		if sesh {
			session := filepath.Base(worktreePath)
			execCmd = exec.Command("tmux", "send-keys", "-t", session, command, "C-m")
			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr
			execCmd.Stdin = os.Stdin

			styledCommandText = greenStyle.Render(fmt.Sprintf("tmux send-keys -t %s ", session)) +
				orangeStyle.Render(command) + greenStyle.Render(" C-m")
		} else {
			execCmd = exec.Command("sh", "-c", command)
			execCmd.Dir = worktreePath
			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr
			execCmd.Stdin = os.Stdin

			styledCommandText = greenStyle.Render(command)
		}

		boldStyle := lipgloss.NewStyle().Bold(true)
		var text string
		if worktree != "" {
			text = boldStyle.Render(fmt.Sprintf("️➡️ Running init command %d of %s in worktree %s: %s...",
				i+1, strconv.Itoa(len(commands)), orangeStyle.Render(worktree), styledCommandText))
		} else {
			text = boldStyle.Render(fmt.Sprintf("️➡️ Running init command %d of %s: %s...",
				i+1, strconv.Itoa(len(commands)), styledCommandText))
		}
		fmt.Println(text)
		if err := execCmd.Run(); err != nil {
			return fmt.Errorf("error running init command '%s': %w", command, err)
		}
	}
	return nil
}
