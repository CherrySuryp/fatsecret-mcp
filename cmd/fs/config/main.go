package main

import (
	"log"

	"github.com/cherrysuryp/fatsecret-mcp/internal/config"
	"github.com/cherrysuryp/fatsecret-mcp/internal/logging"
)

func main() {
	// 1. Init config
	cfg := config.MustLoadConfig()

	// 2. Init logger
	logLevel, err := cfg.FSMCPConfig.SlogLevel()
	if err != nil {
		log.Fatal(err)
	}
	logger := logging.NewLogger(logLevel)

	// 3. Warn user about missing fatsecret user config
	if !cfg.UserConfigExists() {
		logger.Warn("Fatsecret user config file not found, user-specific features will be unavailable")
	} else {
		if cfg.UserConfigEmpty() {
			logger.Warn("Fatsecret user config is empty or invalid, user-specific features will be unavailable")
		}
	}
}
