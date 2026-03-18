package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v1/food.delete_favorite
const foodDeleteFavoritePath = "food/favorite/v1"

// FoodDeleteFavoriteParams holds the parameters for the food.delete_favorite endpoint.
type FoodDeleteFavoriteParams struct {
	FoodID        string  // required
	ServingID     string  // optional
	NumberOfUnits float64 // optional; omitted when zero
}

// FoodDeleteFavorite removes a food from the authenticated user's favorites.
//
// https://platform.fatsecret.com/docs/v1/food.delete_favorite
func (c *FSProfileClient) FoodDeleteFavorite(p FoodDeleteFavoriteParams) (*SuccessResponse, error) {
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

	data, err := c.delete(foodDeleteFavoritePath, params)
	if err != nil {
		return nil, fmt.Errorf("FoodDeleteFavorite: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("FoodDeleteFavorite: decode response: %w", err)
	}

	return &resp, nil
}

