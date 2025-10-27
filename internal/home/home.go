package home

import (
	"os"
	"strings"
)

type Home interface {
	ShortenHome(path string) (string, error)
	ExpandHome(path string) (string, error)
}

type RealHome struct {
}

func NewHome() Home {
	return &RealHome{}
}

func (p *RealHome) ShortenHome(path string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(path, home) {
		return strings.Replace(path, home, "~", 1), nil
	}

	return path, nil
}

func (p *RealHome) ExpandHome(path string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(path, "~") {
		return strings.Replace(path, "~", home, 1), nil
	}
	return path, nil
}
