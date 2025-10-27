package shell

import (
	"github.com/jcelaya775/gwt/internal/home"
	"os/exec"
	"strings"
)

type Shell interface {
	Cmd(cmd string, arg ...string) (string, error)
	ListCmd(cmd string, arg ...string) ([]string, error)
	PrepareCmd(cmd string, replacements map[string]string) ([]string, error)
}

type RealShell struct {
	home home.Home
}

func NewShell(home home.Home) Shell {
	return &RealShell{home}
}

func (c *RealShell) Cmd(cmd string, args ...string) (string, error) {
	foundCmd, err := exec.LookPath(cmd)
	if err != nil {
		return "", err
	}
	output, err := exec.Command(foundCmd, args...).CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func (c *RealShell) ListCmd(cmd string, arg ...string) ([]string, error) {
	command := exec.Command(cmd, arg...)
	output, err := command.Output()
	return strings.Split(string(output), "\n"), err
}

func (c *RealShell) PrepareCmd(cmd string, replacements map[string]string) ([]string, error) {
	cmdParts := strings.Split(cmd, " ")
	result := make([]string, len(cmdParts))

	for i, arg := range cmdParts {
		if strings.HasPrefix(arg, "~") {
			expanded, err := c.home.ExpandHome(arg)
			if err != nil {
				return nil, err
			}
			result[i] = expanded
			continue
		}

		if replacement, ok := replacements[arg]; ok {
			result[i] = replacement
		} else {
			result[i] = arg
		}
	}

	return result, nil
}
