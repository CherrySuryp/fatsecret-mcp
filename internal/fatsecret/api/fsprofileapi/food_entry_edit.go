package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v1/food_entry.edit
const foodEntryEditPath = "food-entries/v1"

// FoodEntryEditParams holds the parameters for the food_entry.edit endpoint.
type FoodEntryEditParams struct {
	FoodEntryID   string // required
	FoodEntryName string // optional
	ServingID     string // optional
	NumberOfUnits string // optional; decimal as string
	Meal          string // optional: "breakfast", "lunch", "dinner", or "other"
}

// FoodEntryEdit updates an existing food diary entry for the authenticated user.
//
// https://platform.fatsecret.com/docs/v1/food_entry.edit
func (c *FSProfileClient) FoodEntryEdit(p FoodEntryEditParams) (*SuccessResponse, error) {
	params := map[string]string{
		"format":        "json",
		"food_entry_id": p.FoodEntryID,
	}
	if p.FoodEntryName != "" {
		params["food_entry_name"] = p.FoodEntryName
	}
	if p.ServingID != "" {
		params["serving_id"] = p.ServingID
	}
	if p.NumberOfUnits != "" {
		params["number_of_units"] = p.NumberOfUnits
	}
	if p.Meal != "" {
		params["meal"] = p.Meal
	}

	data, err := c.put(foodEntryEditPath, params)
	if err != nil {
		return nil, fmt.Errorf("FoodEntryEdit: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("FoodEntryEdit: decode response: %w", err)
	}

	return &resp, nil
}
