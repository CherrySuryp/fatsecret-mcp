package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cherrysuryp/fatsecret-mcp/internal/fatsecret/fsauth"
	"github.com/cherrysuryp/fatsecret-mcp/internal/fatsecret/fsclient"
)

func main() {
	cfg, err := fsauth.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client := fsclient.NewClient(cfg)
	foods, err := client.GetAllFavorites()
	if err != nil {
		log.Fatalf("GetAllFavorites failed: %v", err)
	}

	out, err := json.MarshalIndent(foods, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}

	fmt.Println(string(out))
}
