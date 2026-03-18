package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type FSMCPServerConfig struct {
	LogLevel string `env:"FATSECRET_MCP_SERVER_LOG_LEVEL" envDefault:"INFO"`
}

type FSAPIClientConfig struct {
	ConsumerID     string `env:"FATSECRET_MCP_API_CONSUMER_ID,required"`
	ConsumerSecret string `env:"FATSECRET_MCP_API_CONSUMER_SECRET,required"`
	UserConfigPath string `env:"FATSECRET_MCP_API_USER_CONFIG_PATH" envDefault:"~/.fatsecret-mcp-user-creds.json"`
}

type FSAPIUserConfig struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	SecretToken string `json:"access_token_secret"`
}

type Config struct {
	FSMCPConfig       FSMCPServerConfig
	FSAPIClientConfig FSAPIClientConfig
	FSAPIUserConfig   FSAPIUserConfig
}

/////////////////////////////////
// MCP Server Config methods
/////////////////////////////////

func (c *Config) loadMCPConfig() error {
	if err := env.Parse(&c.FSMCPConfig); err != nil {
		return err
	}
	return nil
}

// SlogLevel parses the LogLevel string into a slog.Level value.
// Returns an error if the value is not one of DEBUG, INFO, WARN, or ERROR.
func (c *FSMCPServerConfig) SlogLevel() (slog.Level, error) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(strings.ToUpper(c.LogLevel))); err != nil {
		return 0, fmt.Errorf("invalid log level %q: must be one of DEBUG, INFO, WARN, ERROR", c.LogLevel)
	}
	return level, nil
}

/////////////////////////////////
// API Client methods
/////////////////////////////////

func (c *Config) loadAPIClientConfig() error {
	if err := env.Parse(&c.FSAPIClientConfig); err != nil {
		return err
	}
	return nil
}

/////////////////////////////////
// API User Client Config methods
/////////////////////////////////

func (c *Config) userConfigPath() (string, error) {
	filePath := c.FSAPIClientConfig.UserConfigPath

	// Handle "~/" path prefix
	if strings.HasPrefix(filePath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		// Replace "~/" with absolute path.
		// For example: ~/.my-config.json -> /Users/johndoe/.my-config.json
		filePath = home + filePath[1:]
	}

	_, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}

	return filePath, nil
}

func (c *Config) loadUserConfig() error {
	filePath, err := c.userConfigPath()
	if err != nil {
		return err
	}

	if filePath == "" {
		return nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	err = json.Unmarshal(data, &c.FSAPIUserConfig)

	return err
}

// UserConfigExists reports whether the user config file exists on disk.
// Errors from stat are treated as absence and return false.
func (c *Config) UserConfigExists() bool {
	filePath, err := c.userConfigPath()
	if err != nil {
		return false
	}
	return filePath != ""
}

// UserConfigEmpty reports whether any required user credential field is blank.
// Should only be called when UserConfigExists returns true.
func (c *Config) UserConfigEmpty() bool {
	if c.FSAPIUserConfig.UserID == "" || c.FSAPIUserConfig.AccessToken == "" || c.FSAPIUserConfig.SecretToken == "" {
		return true
	}
	return false
}

// UserConfigValid reports whether the user config file exists and all credential fields are populated.
// Use this to gate user-specific MCP endpoints.
func (c *Config) UserConfigValid() bool {
	return c.UserConfigExists() && !c.UserConfigEmpty()
}

/////////////////////////////////
// Create New Config
/////////////////////////////////

// MustLoadConfig loads all configuration from environment variables and the user config file.
// Calls log.Fatal on any error — intended for use at process startup only.
func MustLoadConfig() *Config {
	godotenv.Load()

	var cfg Config

	// 1. Load FSMCPConfig
	if err := cfg.loadMCPConfig(); err != nil {
		log.Fatal(err)
	}

	// 2. Load FSAPIClientConfig
	if err := cfg.loadAPIClientConfig(); err != nil {
		log.Fatal(err)
	}

	// 3. Load FSAPIUserConfig
	if err := cfg.loadUserConfig(); err != nil {
		log.Fatal(err)
	}

	return &cfg
}
