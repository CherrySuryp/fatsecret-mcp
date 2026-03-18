package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v1/saved_meal.delete
const savedMealDeletePath = "saved-meals/v1"

// SavedMealDelete removes a saved meal for the authenticated user.
//
// https://platform.fatsecret.com/docs/v1/saved_meal.delete
func (c *FSProfileClient) SavedMealDelete(savedMealID string) (*SuccessResponse, error) {
	params := map[string]string{
		"format":        "json",
		"saved_meal_id": savedMealID,
	}

	data, err := c.delete(savedMealDeletePath, params)
	if err != nil {
		return nil, fmt.Errorf("SavedMealDelete: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("SavedMealDelete: decode response: %w", err)
	}

	return &resp, nil
}

