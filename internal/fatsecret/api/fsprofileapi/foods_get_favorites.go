package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v2/foods.get_favorites
const foodsGetFavoritesPath = "food/favorites/v2"

// FoodsGetFavoritesResponse is the response from the FoodsGetFavorites endpoint.
type FoodsGetFavoritesResponse struct {
	Foods FoodsGetFavoritesResult `json:"foods"`
}

// FoodsGetFavoritesResult holds the list of favorite foods.
type FoodsGetFavoritesResult struct {
	Foods []FavoriteFoodItem `json:"food"`
}

// FavoriteFoodItem is a single favorite food entry.
type FavoriteFoodItem struct {
	FoodID          string `json:"food_id"`
	FoodName        string `json:"food_name"`
	FoodType        string `json:"food_type"`
	FoodURL         string `json:"food_url"`
	FoodDescription string `json:"food_description"`
	ServingID       string `json:"serving_id"`
	NumberOfUnits   string `json:"number_of_units"`
}

// FoodsGetFavorites returns the list of foods saved as favorites for the
// authenticated user.
//
// https://platform.fatsecret.com/docs/v2/foods.get_favorites
func (c *FSProfileClient) FoodsGetFavorites() (*FoodsGetFavoritesResponse, error) {
	params := map[string]string{"format": "json"}

	data, err := c.get(foodsGetFavoritesPath, params)
	if err != nil {
		return nil, fmt.Errorf("FoodsGetFavorites: %w", err)
	}

	var resp FoodsGetFavoritesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("FoodsGetFavorites: decode response: %w", err)
	}

	return &resp, nil
}

