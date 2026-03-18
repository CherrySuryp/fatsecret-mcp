package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v2/foods.get_most_eaten
const foodsGetMostEatenPath = "food/most-eaten/v2"

// FoodsGetMostEatenParams holds the optional parameters for the foods.get_most_eaten endpoint.
type FoodsGetMostEatenParams struct {
	Meal string // optional: "breakfast", "lunch", "dinner", or "other"
}

// FoodsGetMostEatenResponse is the response from the foods.get_most_eaten endpoint.
type FoodsGetMostEatenResponse struct {
	Foods FoodsGetMostEatenResult `json:"foods"`
}

// FoodsGetMostEatenResult holds the list of most-eaten foods.
type FoodsGetMostEatenResult struct {
	Foods []MostEatenFoodItem `json:"food"`
}

// MostEatenFoodItem is a single most-eaten food entry.
type MostEatenFoodItem struct {
	FoodID        string `json:"food_id"`
	FoodName      string `json:"food_name"`
	FoodType      string `json:"food_type"`
	FoodURL       string `json:"food_url"`
	ServingID     string `json:"serving_id"`
	NumberOfUnits string `json:"number_of_units"`
}

// FoodsGetMostEaten returns the foods most frequently eaten by the authenticated user,
// optionally filtered by meal type.
//
// https://platform.fatsecret.com/docs/v2/foods.get_most_eaten
func (c *FSProfileClient) FoodsGetMostEaten(p FoodsGetMostEatenParams) (*FoodsGetMostEatenResponse, error) {
	params := map[string]string{"format": "json"}
	if p.Meal != "" {
		params["meal"] = p.Meal
	}

	data, err := c.get(foodsGetMostEatenPath, params)
	if err != nil {
		return nil, fmt.Errorf("FoodsGetMostEaten: %w", err)
	}

	var resp FoodsGetMostEatenResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("FoodsGetMostEaten: decode response: %w", err)
	}

	return &resp, nil
}

