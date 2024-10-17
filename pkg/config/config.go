package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const Version = "v1"

type SlackConfig struct {
	Token   string `yaml:"token"`   // Slack token
	Channel string `yaml:"channel"` // Slack channel
}

type Config struct {
	Slack SlackConfig `yaml:"slack"`
}

// ConfigFilePath defines the path where the configuration file will be stored
// ~/.config/gcp-pam-jit-v1.yaml
var path = filepath.Join(os.Getenv("HOME"), ".config", fmt.Sprintf("gcp-pam-jit-%s.yaml", Version))

// LoadConfig loads the configuration from the YAML file
func LoadConfig() (*Config, error) {
	// Check if the config file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found at %s", path)
	}

	// Read the file contents
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing YAML: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to the YAML file
func SaveConfig(config *Config) error {

	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling YAML: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}
