package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v1/saved_meal.edit
const savedMealEditPath = "saved-meals/v1"

// SavedMealEditParams holds the parameters for the saved_meal.edit endpoint.
type SavedMealEditParams struct {
	SavedMealID          string // required
	SavedMealName        string // optional
	SavedMealDescription string // optional
	Meals                string // optional; comma-separated meal types (e.g. "breakfast,lunch")
}

// SavedMealEdit updates the name, description, or meal assignments of an existing saved meal.
//
// https://platform.fatsecret.com/docs/v1/saved_meal.edit
func (c *FSProfileClient) SavedMealEdit(p SavedMealEditParams) (*SuccessResponse, error) {
	params := map[string]string{
		"format":        "json",
		"saved_meal_id": p.SavedMealID,
	}
	if p.SavedMealName != "" {
		params["saved_meal_name"] = p.SavedMealName
	}
	if p.SavedMealDescription != "" {
		params["saved_meal_description"] = p.SavedMealDescription
	}
	if p.Meals != "" {
		params["meals"] = p.Meals
	}

	data, err := c.put(savedMealEditPath, params)
	if err != nil {
		return nil, fmt.Errorf("SavedMealEdit: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("SavedMealEdit: decode response: %w", err)
	}

	return &resp, nil
}

