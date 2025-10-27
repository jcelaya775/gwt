package connector

import (
	"errors"
	"fmt"
	"github.com/jcelaya775/gwt/internal/shell"
	"os/exec"
)

type Connector struct {
	shell shell.Shell
}

func New(shell shell.Shell) *Connector {
	return &Connector{shell: shell}
}

func (c *Connector) SeshConnect(dir string) error {
	_, err := c.shell.Cmd("sesh", "connect", dir)
	return err
}

func (c *Connector) WebstormConnect(dir string) error {
	return c.jetbrainsConnect(WebStorm, dir)
}

func (c *Connector) IntelliJConnect(dir string) error {
	return c.jetbrainsConnect(IntelliJIDEA, dir)
}

func (c *Connector) PyCharmConnect(dir string) error {
	return c.jetbrainsConnect(PyCharm, dir)
}

func (c *Connector) CLionConnect(dir string) error {
	return c.jetbrainsConnect(CLion, dir)
}

func (c *Connector) RiderConnect(dir string) error {
	return c.jetbrainsConnect(Rider, dir)
}

func (c *Connector) GoLandConnect(dir string) error {
	return c.jetbrainsConnect(GoLand, dir)
}

func (c *Connector) DataGripConnect(dir string) error {
	return c.jetbrainsConnect(DataGrip, dir)
}

type jetbrainsIDE string

const (
	WebStorm     jetbrainsIDE = "webstorm"
	IntelliJIDEA jetbrainsIDE = "idea"
	PyCharm      jetbrainsIDE = "pycharm"
	CLion        jetbrainsIDE = "clion"
	Rider        jetbrainsIDE = "rider"
	GoLand       jetbrainsIDE = "goland"
	DataGrip     jetbrainsIDE = "datagrip"
)

func (c *Connector) jetbrainsConnect(ide jetbrainsIDE, dir string) error {
	if _, err := exec.LookPath(string(ide)); err != nil {
		return fmt.Errorf("%s command not found in PATH", ide)
	}

	var err error
	if _, err = exec.LookPath("setsid"); err == nil {
		cmd := fmt.Sprintf("%s %s >/dev/null 2>&1 &", ide, dir)
		output, cmdErr := c.shell.Cmd("setsid", "sh", "-c", cmd)
		if cmdErr != nil {
			err = errors.New(output)
		}
	} else if _, err = exec.LookPath("nohup"); err == nil {
		cmd := fmt.Sprintf("%s %s >/dev/null 2>&1 &", ide, dir)
		output, cmdErr := c.shell.Cmd("nohup", cmd)
		if cmdErr != nil {
			err = errors.New(output)
		}
	} else if _, err = exec.LookPath("open"); err == nil {
		output, cmdErr := c.shell.Cmd("open", "-a", "WebStorm", dir)
		if cmdErr != nil {
			err = errors.New(output)
		}
	} else if _, err = exec.LookPath("systemd-run"); err == nil {
		cmd := fmt.Sprintf("%s %s >/dev/null 2>&1 &", ide, dir)
		output, cmdErr := c.shell.Cmd("systemd-run", "--user", "--scope", "sh", "-c", cmd)
		if cmdErr != nil {
			err = errors.New(output)
		}
	} else {
		output, cmdErr := c.shell.Cmd("%s", string(ide), fmt.Sprintf("%s /dev/null 2>&1 &", dir))
		if cmdErr != nil {
			err = errors.New(output)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to launch %s: %v", ide, err)
	}
	return nil
}

//VsCodeConnect(c *AddConfig)
//CursorConnect(c *AddConfig)
//GetZoxidePath(c *AddConfig) string
