package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration file
type Config struct {
	Tests []TestConfig `yaml:"tests"`
}

// TestConfig represents a test configuration
type TestConfig struct {
	Name     string `yaml:"name"`
	Command  string `yaml:"command"`
	Blocking bool   `yaml:"blocking"`
	Timeout  int    `yaml:"timeout"` // in seconds
}

// Load loads the configuration from a file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	for i := range config.Tests {
		if config.Tests[i].Timeout == 0 {
			config.Tests[i].Timeout = 300 // 5 minutes default
		}
	}

	return &config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if len(c.Tests) == 0 {
		return fmt.Errorf("no tests configured")
	}

	for _, test := range c.Tests {
		if test.Name == "" {
			return fmt.Errorf("test name cannot be empty")
		}
		if test.Command == "" {
			return fmt.Errorf("test command cannot be empty for test '%s'", test.Name)
		}
		if test.Timeout < 0 {
			return fmt.Errorf("test timeout must be positive for test '%s'", test.Name)
		}
	}

	return nil
}
