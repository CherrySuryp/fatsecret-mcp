package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v1/food.add_favorite
const foodAddFavoritePath = "food/favorite/v1"

// FoodAddFavoriteParams holds the parameters for the food.add_favorite endpoint.
type FoodAddFavoriteParams struct {
	FoodID        string  // required
	ServingID     string  // optional
	NumberOfUnits float64 // optional; omitted when zero
}

// FoodAddFavorite adds a food to the authenticated user's favorites.
//
// https://platform.fatsecret.com/docs/v1/food.add_favorite
func (c *FSProfileClient) FoodAddFavorite(p FoodAddFavoriteParams) (*SuccessResponse, error) {
	params := map[string]string{
		"format":  "json",
		"food_id": p.FoodID,
	}
	if p.ServingID != "" {
		params["serving_id"] = p.ServingID
	}
	if p.NumberOfUnits != 0 {
		params["number_of_units"] = fmt.Sprintf("%g", p.NumberOfUnits)
	}

	data, err := c.post(foodAddFavoritePath, params)
	if err != nil {
		return nil, fmt.Errorf("FoodAddFavorite: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("FoodAddFavorite: decode response: %w", err)
	}

	return &resp, nil
}

