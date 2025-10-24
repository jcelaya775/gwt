package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Version         string   `yaml:"version"`
	InitCommands    []string `yaml:"init_commands,omitempty"` // Commands to run when initializing a repository
	Defaults        Defaults `yaml:"defaults,omitempty"`
	DestroyCommands []string `yaml:"destroy_commands,omitempty"` // Commands to run when destroying a worktree
}

type Defaults struct {
	BaseBranch string `yaml:"base_branch,omitempty"` // Default base branch for new worktrees
}

const (
	ConfigFileName    = ".gwt.yml"
	CurrentVersion    = "1.0"
	DefaultBaseBranch = "main"
)

// LoadConfig loads configuration from .gwt.yml in the repository root
func LoadConfig(repoRoot string) (*Config, error) {
	configPath := filepath.Join(repoRoot, ConfigFileName)

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			Version: CurrentVersion,
			Defaults: Defaults{
				BaseBranch: DefaultBaseBranch,
			},
		}, nil
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if c.Version == "" {
		c.Version = CurrentVersion
	}

	if c.Defaults.BaseBranch == "" {
		c.Defaults.BaseBranch = DefaultBaseBranch
	}

	return nil
}
