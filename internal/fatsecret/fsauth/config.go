package fsauth

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ClientID          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	AccessToken       string `json:"accessToken,omitempty"`
	AccessTokenSecret string `json:"accessTokenSecret,omitempty"`
	UserID            string `json:"userId,omitempty"`
}

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".fatsecret-mcp-config.json")
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}

	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		return nil, err
	}

	// Env vars take precedence for credentials
	if cfg.ClientID == "" {
		cfg.ClientID = loaded.ClientID
	}
	if cfg.ClientSecret == "" {
		cfg.ClientSecret = loaded.ClientSecret
	}
	cfg.AccessToken = loaded.AccessToken
	cfg.AccessTokenSecret = loaded.AccessTokenSecret
	cfg.UserID = loaded.UserID

	return cfg, nil
}

func SaveConfig(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(), data, 0600)
}
