package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/cherrysuryp/fatsecret-mcp/internal/fatsecret/fsauth"
	"github.com/cherrysuryp/fatsecret-mcp/internal/fatsecret/fsclient"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: addfavorite <food_id> [serving_id] [number_of_units]")
	}

	var req fsclient.FoodAddFavoriteReq
	req.FoodID = os.Args[1]

	if len(os.Args) > 2 {
		req.ServingID = &os.Args[2]
	}

	if len(os.Args) > 3 {
		n, err := strconv.ParseUint(os.Args[3], 10, 64)
		if err != nil {
			log.Fatalf("invalid number_of_units: %v", err)
		}
		nUint := uint(n)
		req.NumberOfUnits = &nUint
	}

	cfg, err := fsauth.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client := fsclient.NewClient(cfg)
	resp, err := client.FoodAddFavorite(req)
	if err != nil {
		log.Fatalf("FoodAddFavorite failed: %v", err)
	}

	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}

	fmt.Println(string(out))
}
