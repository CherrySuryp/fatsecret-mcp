package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v2/foods.get_recently_eaten
const foodsGetRecentlyEatenPath = "food/recently-eaten/v2"

// FoodsGetRecentlyEatenParams holds the optional parameters for the foods.get_recently_eaten endpoint.
type FoodsGetRecentlyEatenParams struct {
	Meal string // optional: "breakfast", "lunch", "dinner", or "other"
}

// FoodsGetRecentlyEatenResponse is the response from the foods.get_recently_eaten endpoint.
type FoodsGetRecentlyEatenResponse struct {
	Foods FoodsGetRecentlyEatenResult `json:"foods"`
}

// FoodsGetRecentlyEatenResult holds the list of recently eaten foods.
type FoodsGetRecentlyEatenResult struct {
	Foods []RecentlyEatenFoodItem `json:"food"`
}

// RecentlyEatenFoodItem is a single recently eaten food entry.
type RecentlyEatenFoodItem struct {
	FoodID        string `json:"food_id"`
	FoodName      string `json:"food_name"`
	FoodType      string `json:"food_type"`
	FoodURL       string `json:"food_url"`
	ServingID     string `json:"serving_id"`
	NumberOfUnits string `json:"number_of_units"`
}

// FoodsGetRecentlyEaten returns the foods most recently eaten by the authenticated user,
// optionally filtered by meal type.
//
// https://platform.fatsecret.com/docs/v2/foods.get_recently_eaten
func (c *FSProfileClient) FoodsGetRecentlyEaten(p FoodsGetRecentlyEatenParams) (*FoodsGetRecentlyEatenResponse, error) {
	params := map[string]string{"format": "json"}
	if p.Meal != "" {
		params["meal"] = p.Meal
	}

	data, err := c.get(foodsGetRecentlyEatenPath, params)
	if err != nil {
		return nil, fmt.Errorf("FoodsGetRecentlyEaten: %w", err)
	}

	var resp FoodsGetRecentlyEatenResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("FoodsGetRecentlyEaten: decode response: %w", err)
	}

	return &resp, nil
}

