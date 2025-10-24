package selecter

import (
	"bytes"
	"github.com/charmbracelet/huh"
	"os"
	"os/exec"
	"strings"
)

type Selecter string

const (
	Fzf Selecter = "fzf"
)

type Select struct {
	Type Selecter
}

func New() *Select {
	if _, err := exec.LookPath("fzf"); err == nil {
		return &Select{
			Type: Fzf,
		}
	}
	return &Select{
		Type: "",
	}
}

func (s *Select) Select(header string, options []string) (string, error) {
	if len(options) == 0 {
		return "", nil
	}

	switch s.Type {
	case Fzf:
		selectedValue, err := s.fzfSelect(header, options)
		if err != nil {
			return "", err
		}
		return selectedValue, nil
	default:
		selectedValue, err := s.defaultSelect(header, options)
		if err != nil {
			return "", err
		}
		return selectedValue, nil
	}
}

func (s *Select) fzfSelect(header string, options []string) (string, error) {
	cmd := exec.Command("fzf", "--header", header)
	cmd.Stdin = strings.NewReader(strings.Join(options, "\n"))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if err.Error() == "exit status 130" {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func (s *Select) defaultSelect(header string, options []string) (string, error) {
	huhOptions := make([]huh.Option[string], len(options))
	for i, option := range options {
		huhOptions[i] = huh.NewOption(option, option)
	}
	var selectedValue string
	err := huh.NewSelect[string]().
		Title(header).
		Options(huhOptions...).
		Value(&selectedValue).
		Run()
	if err != nil {
		return "", err
	}
	return selectedValue, nil
}

func (s *Select) MultiSelect(header string, options []string) ([]string, error) {
	if len(options) == 0 {
		return []string{}, nil
	}

	switch s.Type {
	case Fzf:
		selectedValues, err := s.fzfMultiSelect(header, options)
		if err != nil {
			return nil, err
		}
		return selectedValues, nil
	default:
		selectedValues, err := s.defaultMultiSelect(header, options)
		if err != nil {
			return nil, err
		}
		return selectedValues, nil
	}
}

func (s *Select) fzfMultiSelect(header string, options []string) ([]string, error) {
	cmd := exec.Command("fzf", "--header", header, "--multi")
	cmd.Stdin = strings.NewReader(strings.Join(options, "\n"))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if err.Error() == "exit status 130" {
			return nil, nil
		}
		return nil, err
	}
	selected := strings.Split(strings.TrimSpace(out.String()), "\n")
	return selected, nil
}

func (s *Select) defaultMultiSelect(header string, options []string) ([]string, error) {
	huhOptions := make([]huh.Option[string], len(options))
	for i, option := range options {
		huhOptions[i] = huh.NewOption(option, option)
	}
	var selectedValues []string
	err := huh.NewMultiSelect[string]().
		Title(header).
		Options(huhOptions...).
		Value(&selectedValues).
		Run()
	if err != nil {
		return nil, err
	}
	return selectedValues, nil
}
