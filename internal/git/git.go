package git

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/huh/spinner"
	"github.com/jcelaya775/gwt/internal/config"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Git struct {
	worktreeRoot string
}

func New() (*Git, error) {
	return &Git{}, nil
}

func (g *Git) SetWorktreeRoot() error {
	worktreeRoot, err := getWorktreeRoot()
	g.worktreeRoot = worktreeRoot
	return err
}

func (g *Git) GetWorktreeRoot() string {
	return g.worktreeRoot
}

func (g *Git) GetRepoName() string {
	return filepath.Base(g.worktreeRoot)
}

func (g *Git) CloneRepo(repoURL string, dir string) error {
	var cmdArgs []string
	if dir != "" {
		cmdArgs = []string{"clone", "--no-checkout", repoURL, dir}
	} else {
		cmdArgs = []string{"clone", "--no-checkout", repoURL}
	}
	output, err := exec.Command("git", cmdArgs...).CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		return errors.New(string(output))
	}

	var repoDir string
	if dir != "" {
		repoDir = dir
	} else {
		parts := strings.Split(repoURL, "/")
		repoDir = strings.TrimSuffix(parts[len(parts)-1], ".git")
	}
	repoPath, err := filepath.Abs(repoDir)
	if err != nil {
		return err
	}
	fmt.Println("Repository cloned to:", repoPath)
	g.worktreeRoot = repoPath

	// TODO: Checkout a dummy branch and create worktree for main branch relative to worktree dir (config)
	// TODO: Create pseudo-random branch name to avoid conflicts, random 6 digit number

	output, err = exec.Command("git", "-C", g.worktreeRoot, "branch", "--show-current").CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}
	originalBranch := strings.TrimSpace(string(output))
	fmt.Println("Original branch:", originalBranch)

	dummyBranch := "dummy"
	fmt.Println("Creating dummy branch:", dummyBranch)
	output, err = exec.Command("git", "-C", g.worktreeRoot, "checkout", "-b", dummyBranch).CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}

	var worktreeArgs []string
	worktreeArgs = []string{"-C", g.worktreeRoot, "worktree", "add", originalBranch, originalBranch}
	fmt.Println("Creating worktree for original branch:", originalBranch)
	output, err = exec.Command("git", worktreeArgs...).CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}

	return nil
}

func (g *Git) AddWorktree(config *config.Config, branch string, commitish string, noPull bool, force bool) (string, error) {
	cmdArgs := []string{"-C", g.worktreeRoot, "worktree", "add"}

	var baseBranch string

	existsLocally, err := g.BranchExistsLocally(branch)
	if err != nil {
		return "", err
	}
	existsRemotely, err := g.BranchExistsRemotely(branch)
	if err != nil {
		return "", err
	}

	if !existsRemotely && strings.HasPrefix(branch, "origin/") {
		return "", fmt.Errorf("branch '%s' does not exist remotely. Remove the 'origin/' prefix to create a new branch", branch)
	}

	parsedBranch := strings.TrimPrefix(branch, "origin/")
	if commitish != "" {
		baseBranch = commitish
		cmdArgs = append(cmdArgs, "-b", parsedBranch, parsedBranch, "--checkout", commitish)
	} else if existsLocally {
		baseBranch = parsedBranch
		cmdArgs = append(cmdArgs, parsedBranch, "--checkout", baseBranch)
	} else if existsRemotely {
		baseBranch = branch
		cmdArgs = append(cmdArgs, "-b", parsedBranch, parsedBranch, "--checkout", branch)
	} else {
		baseBranch = config.Defaults.BaseBranch
		cmdArgs = append(cmdArgs, "-b", parsedBranch, parsedBranch, baseBranch)
	}

	worktreeExists, err := g.WorktreeExists(baseBranch)

	if !noPull && !strings.HasPrefix(baseBranch, "origin/") {
		// TODO: Check for uncommitted changes or merge conflicts, and prompt user with confirmation message before pulling
		var err error
		baseBranchPath := filepath.Join(g.worktreeRoot, baseBranch)
		_ = spinner.New().
			Title(fmt.Sprintf("Pulling base branch '%s'... (press ctrl-c to skip)", baseBranch)).
			Action(func() {
				var output []byte
				var innerErr error
				if worktreeExists {
					output, innerErr = exec.Command("git", "-C", baseBranchPath, "pull", "origin", fmt.Sprintf("%s:%s", baseBranch, baseBranch)).CombinedOutput()
				} else if existsLocally && existsRemotely {
					output, innerErr = exec.Command("git", "-C", g.worktreeRoot, "fetch", "origin", fmt.Sprintf("%s:%s", baseBranch, baseBranch)).CombinedOutput()
				}

				if innerErr != nil {
					err = errors.New(string(output))
				}
			}).
			Run()
		if err != nil {
			return "", errors.Join(err, fmt.Errorf("failed to pull base branch. You can retry without pulling using the --no-pull flag"))
		}
	}

	if force {
		cmdArgs = append(cmdArgs, "--force")
	}

	output, err := exec.Command("git", cmdArgs...).CombinedOutput()
	if err != nil {
		return "", errors.New(string(output))
	}
	fmt.Println(string(output))

	worktreePathRegex, err := regexp.Compile("Preparing worktree.*'(.*)'")
	if err != nil {
		return "", err
	}
	matches := worktreePathRegex.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return "", errors.New("could not parse worktree path from git output")
	}
	worktreePath := filepath.Join(g.worktreeRoot, matches[1])
	return worktreePath, nil
}

func (g *Git) RemoveWorktree(worktree string, force, keepBranch bool) error {
	branch, err := g.GetWorktreeBranch(worktree)

	cmdArgs := []string{"-C", g.worktreeRoot, "worktree", "remove"}
	if force {
		cmdArgs = append(cmdArgs, "--force")
	}
	cmdArgs = append(cmdArgs, worktree)

	output, err := exec.Command("git", cmdArgs...).CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}

	if !keepBranch {
		if err != nil {
			return err
		}
		branchOutput, branchErr := exec.Command("git", "-C", g.worktreeRoot, "branch", "-D", branch).CombinedOutput()
		if branchErr != nil {
			return errors.New(string(branchOutput))
		}
	}

	if err = g.removeEmptyParentDirs(worktree); err != nil {
		return err
	}

	return nil
}

func (g *Git) removeEmptyParentDirs(worktree string) error {
	worktreePath := filepath.Join(g.worktreeRoot, worktree)
	parentDir := filepath.Dir(worktreePath)
	for parentDir != g.worktreeRoot {
		dirEntries, err := os.ReadDir(parentDir)
		if err != nil {
			return err
		}
		if len(dirEntries) == 0 {
			if err := os.Remove(parentDir); err != nil {
				return err
			}
			parentDir = filepath.Dir(parentDir)
		} else {
			break
		}
	}
	return nil
}

func (g *Git) ListWorktrees() ([]string, error) {
	output, err := exec.Command("git", "worktree", "list").Output()
	if err != nil {
		return nil, errors.New(string(output))
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	worktrees := make([]string, 0, len(lines))
	for _, line := range lines[1:] {
		worktreeAbsPath := strings.Fields(line)[0]
		worktreeRelPath := strings.TrimPrefix(worktreeAbsPath, g.worktreeRoot)
		worktrees = append(worktrees, worktreeRelPath)
	}
	return worktrees, nil
}

func (g *Git) ListBranches(onlyLocal, hideBranchesWithWorktrees bool) ([]string, error) {
	worktreeBranches := make(map[string]bool)
	if hideBranchesWithWorktrees {
		output, err := exec.Command("git", "worktree", "list").Output()
		if err != nil {
			return nil, errors.New(string(output))
		}

		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range lines {
			parts := strings.Fields(line)
			if len(parts) < 3 {
				return nil, errors.New("could not parse worktree list output")
			}
			worktreeBranch := strings.TrimSuffix(strings.TrimPrefix(parts[2], "["), "]")
			worktreeBranches[worktreeBranch] = true
		}
	}

	ignoreBranches := map[string]bool{
		"origin/HEAD": true,
	}

	cmdArgs := []string{"branch", "--format=%(refname:short)"}
	if !onlyLocal {
		cmdArgs = append(cmdArgs, "-a")
	}
	output, err := exec.Command("git", cmdArgs...).Output()
	if err != nil {
		return nil, errors.New(string(output))
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	branches := make([]string, 0, len(lines))
	for _, line := range lines {
		branch := strings.TrimSpace(line)
		parsedBranch := strings.TrimPrefix(strings.TrimPrefix(branch, "origin"), "/")
		if parsedBranch == "" {
			continue
		}

		if strings.HasPrefix(branch, "origin/") {
			branchWithoutWTExistsLocally, err := g.BranchExistsLocally(parsedBranch)
			if err != nil {
				return nil, err
			}
			if branchWithoutWTExistsLocally {
				continue
			}
		}
		if _, ignore := ignoreBranches[branch]; ignore {
			continue
		}
		if _, branchHasWorktree := worktreeBranches[parsedBranch]; hideBranchesWithWorktrees && branchHasWorktree {
			continue
		}
		if onlyLocal && strings.HasPrefix(branch, "origin/") {
			continue
		}
		branches = append(branches, branch)
	}

	return branches, nil
}

func (g *Git) GetWorktreeBranch(worktree string) (string, error) {
	cmdArgs := []string{"-C", g.worktreeRoot, "worktree", "list"}
	output, err := exec.Command("git", cmdArgs...).Output()
	if err != nil {
		return "", errors.New(string(output))
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		worktreeAbsPath := fields[0]
		worktreeRelPath := strings.TrimPrefix(worktreeAbsPath, g.worktreeRoot)
		if worktreeRelPath == worktree {
			if len(fields) < 3 {
				return "", errors.New("could not determine branch for worktree")
			}
			unparsedBranch := fields[2]
			return strings.TrimSuffix(strings.TrimPrefix(unparsedBranch, "["), "]"), nil
		}
	}
	return "", errors.New("worktree not found")
}

func (*Git) Fetch() error {
	var err error
	_ = spinner.New().
		Title("Syncing with remote... (press ctrl-c to skip)").
		Action(func() {
			cmd := exec.Command("git", "fetch")
			output, innerErr := cmd.CombinedOutput()
			if innerErr != nil {
				err = errors.New(string(output))
			}
		}).
		Run()
	return err
}

func (g *Git) WorktreeExists(worktree string) (bool, error) {
	worktreeBranch, err := g.GetWorktreeBranch(worktree)
	if err != nil {
		if strings.Contains(err.Error(), "worktree not found") {
			return false, nil
		}
		return false, err
	}
	if worktreeBranch != "" {
		return true, nil
	}
	return false, nil
}

func (*Git) BranchExistsLocally(branch string) (bool, error) {
	output, err := exec.Command("git", "branch", "--list", branch).CombinedOutput()
	if err != nil {
		return false, errors.New(string(output))
	}

	if strings.TrimSpace(string(output)) == branch {
		return true, nil
	}
	return false, nil
}

func (*Git) BranchExistsRemotely(branch string) (bool, error) {
	if !strings.HasPrefix(branch, "origin/") {
		branch = "origin/" + branch
	}
	output, err := exec.Command("git", "branch", "-r", "--list", branch).CombinedOutput()
	if err != nil {
		return false, errors.New(string(output))
	}

	if strings.TrimSpace(string(output)) == branch {
		return true, nil
	}
	return false, nil
}

func getWorktreeRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-common-dir")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.New("not in a git repository")
	}

	gitDir := strings.TrimSpace(string(output))
	if strings.HasSuffix(gitDir, ".git") {
		gitDirAbsolute, err := filepath.Abs(gitDir)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s/", filepath.Dir(gitDirAbsolute)), nil
	} else {
		return "", errors.New("could not find git repository root containing .git. Please use gwt clone to clone the repository")
	}
}
