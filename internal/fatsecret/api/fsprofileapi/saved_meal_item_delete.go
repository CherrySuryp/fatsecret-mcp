package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v1/saved_meal_item.delete
const savedMealItemDeletePath = "saved-meals/item/v1"

// SavedMealItemDelete removes an item from a saved meal.
//
// https://platform.fatsecret.com/docs/v1/saved_meal_item.delete
func (c *FSProfileClient) SavedMealItemDelete(savedMealItemID string) (*SuccessResponse, error) {
	params := map[string]string{
		"format":             "json",
		"saved_meal_item_id": savedMealItemID,
	}

	data, err := c.delete(savedMealItemDeletePath, params)
	if err != nil {
		return nil, fmt.Errorf("SavedMealItemDelete: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("SavedMealItemDelete: decode response: %w", err)
	}

	return &resp, nil
}

