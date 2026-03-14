package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cherrysuryp/fatsecret-mcp/internal/fatsecret/fsauth"
	"github.com/cherrysuryp/fatsecret-mcp/internal/fatsecret/fsclient"
)

func main() {
	cfg, err := fsauth.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	var meal *fsclient.MealType
	if len(os.Args) > 1 {
		m := fsclient.MealType(os.Args[1])
		meal = &m
	}

	client := fsclient.NewClient(cfg)
	foods, err := client.GetMostEaten(fsclient.GetMostEatenReq{Meal: meal})
	if err != nil {
		log.Fatalf("GetMostEaten failed: %v", err)
	}

	out, err := json.MarshalIndent(foods, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}

	fmt.Println(string(out))
}
