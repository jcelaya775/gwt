package git

type AddWorktreeOptions struct {
	Branch    string
	Commitish string
	Pull      bool
	Force     bool
}
