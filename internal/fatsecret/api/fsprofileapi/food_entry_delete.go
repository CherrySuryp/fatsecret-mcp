package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v1/food_entry.delete
const foodEntryDeletePath = "food-entries/v1"

// FoodEntryDelete removes a food diary entry for the authenticated user.
//
// https://platform.fatsecret.com/docs/v1/food_entry.delete
func (c *FSProfileClient) FoodEntryDelete(foodEntryID string) (*SuccessResponse, error) {
	params := map[string]string{
		"format":        "json",
		"food_entry_id": foodEntryID,
	}

	data, err := c.delete(foodEntryDeletePath, params)
	if err != nil {
		return nil, fmt.Errorf("FoodEntryDelete: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("FoodEntryDelete: decode response: %w", err)
	}

	return &resp, nil
}
