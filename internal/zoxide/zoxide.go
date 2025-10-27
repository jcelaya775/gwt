package zoxide

import (
	"github.com/jcelaya775/gwt/internal/shell"
)

type Zoxide struct {
	shell shell.Shell
}

func New(shell shell.Shell) *Zoxide {
	return &Zoxide{shell: shell}
}

func (z *Zoxide) AddPath(path string) error {
	_, err := z.shell.Cmd("zoxide", "add", path)
	return err
}

func (z *Zoxide) RemovePath(path string) error {
	_, err := z.shell.Cmd("zoxide", "remove", path)
	return err
}

//func (z *Zoxide) AddZoxideEntries(c *config.Config) {
//	baseName := c.GetZoxideBasePath()
//
//	var foldersToAdd []string
//	foldersToAdd = append(foldersToAdd, baseName)
//
//	foldersToAdd = addConfigFolders(foldersToAdd, c.ZoxideFolders, baseName, c.DirectoryReader)
//
//	for _, folder := range foldersToAdd {
//		err := z.AddPath(folder)
//		util.CheckError(err)
//	}
//}
