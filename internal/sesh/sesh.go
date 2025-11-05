package sesh

import (
	"path"
	"strings"
	"unicode"

	"github.com/jcelaya775/gwt/internal/shell"
)

type Sesh struct {
	shell shell.Shell
}

func New(shell shell.Shell) *Sesh {
	return &Sesh{shell: shell}
}

func (s *Sesh) SessionExists(worktree string) (bool, error) {
	sessionName := normalize(path.Base(worktree))

	sessions, err := s.shell.Cmd("sesh", "list")
	if err != nil {
		return false, err
	}

	for _, session := range strings.Split(sessions, "\n") {
		session = strings.TrimSpace(session)
		if session == "" {
			continue
		}

		if normalize(session) == sessionName {
			return true, nil
		}
	}

	return false, nil
}

// normalize makes a name comparable by:
// - trimming and lowercasing
// - removing runes that are not letters, digits, '-' or '_'
// This effectively strips emojis and other symbols so "my-branch" matches "my-branch 🤖".
func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
