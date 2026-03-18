package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v1/saved_meal_item.edit
const savedMealItemEditPath = "saved-meals/item/v1"

// SavedMealItemEditParams holds the parameters for the saved_meal_item.edit endpoint.
type SavedMealItemEditParams struct {
	SavedMealItemID   string // required
	SavedMealItemName string // optional
	NumberOfUnits     string // optional
}

// SavedMealItemEdit updates the name or serving amount of an existing saved meal item.
//
// https://platform.fatsecret.com/docs/v1/saved_meal_item.edit
func (c *FSProfileClient) SavedMealItemEdit(p SavedMealItemEditParams) (*SuccessResponse, error) {
	params := map[string]string{
		"format":             "json",
		"saved_meal_item_id": p.SavedMealItemID,
	}
	if p.SavedMealItemName != "" {
		params["saved_meal_item_name"] = p.SavedMealItemName
	}
	if p.NumberOfUnits != "" {
		params["number_of_units"] = p.NumberOfUnits
	}

	data, err := c.put(savedMealItemEditPath, params)
	if err != nil {
		return nil, fmt.Errorf("SavedMealItemEdit: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("SavedMealItemEdit: decode response: %w", err)
	}

	return &resp, nil
}

