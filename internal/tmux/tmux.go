package tmux

import (
	"errors"
	"fmt"
	"github.com/jcelaya775/gwt/internal/shell"
	"os/exec"
	"strings"
	"time"
)

type Tmux struct {
	shell *shell.Shell
}

func NewTmux() *Tmux {
	return &Tmux{}
}

func (t *Tmux) SendKeys(session string, cmdStr string) error {
	tmuxCmd := exec.Command("tmux", "send-keys", "-t", session, cmdStr, "C-m")
	tmuxCmd.Stdout = nil
	tmuxCmd.Stderr = nil
	if err := tmuxCmd.Run(); err != nil {
		return fmt.Errorf("error sending command to tmux session '%s': %w", session, err)
	}
	return nil
}

// waitForTmuxWindowActive waits until the tmux target session:window reports window_active == 1.
// `window` may be a window index ("0") or name ("editor"). Timeout controls how long to wait.
func waitForTmuxWindowActive(session, window string, timeout time.Duration) error {
	target := fmt.Sprintf("%s:%s", session, window)
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		out, err := exec.Command("tmux", "display-message", "-p", "-t", target, "#{window_active}").Output()
		if err == nil {
			if strings.TrimSpace(string(out)) == "1" {
				return nil
			}
		}
		// short sleep then retry
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("timed out waiting for tmux target %s to become active", target)
}

// waitForTmuxSessionReady waits until the tmux session exists and has an active window.
// Use this when you don't know the exact window name/index but want the session to be ready.
func waitForTmuxSessionReady(session string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		// check session existence
		if err := exec.Command("tmux", "has-session", "-t", session).Run(); err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// list windows and check active flag
		out, err := exec.Command("tmux", "list-windows", "-t", session, "-F", "#{window_index} #{window_name} #{window_active}").Output()
		if err == nil {
			for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
				if line == "" {
					continue
				}
				fields := strings.Fields(line)
				if len(fields) >= 3 && fields[len(fields)-1] == "1" {
					return nil // found active window
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("timed out waiting for tmux session to be ready")
}

// Example usage inside your loop before sending the keys:
//
// session := filepath.Base(worktreePath)
// // if you know a target window name or index, prefer waitForTmuxWindowActive(session, "<win>", 5*time.Second)
// if err := waitForTmuxSessionReady(session, 5*time.Second); err != nil {
//     // fallback or return error
// }
// // now send keys
// tmuxCmd := exec.Command("tmux", "send-keys", "-t", session, cmdStr, "C-m")
// tmuxCmd.Stdout = os.Stdout
// tmuxCmd.Stderr = os.Stderr
// if err := tmuxCmd.Run(); err != nil {
//     return fmt.Errorf("error sending init command to tmux session '%s': %w", session, err)
// }
