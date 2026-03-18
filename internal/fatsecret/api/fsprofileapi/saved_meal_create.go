package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v1/saved_meal.create
const savedMealCreatePath = "saved-meals/v1"

// SavedMealCreateResponse is the response from the saved_meal.create endpoint.
type SavedMealCreateResponse struct {
	SavedMealID struct {
		Value string `json:"value"`
	} `json:"saved_meal_id"`
}

// SavedMealCreateParams holds the parameters for the saved_meal.create endpoint.
type SavedMealCreateParams struct {
	SavedMealName        string // required
	SavedMealDescription string // optional
	Meals                string // optional; comma-separated meal types (e.g. "breakfast,lunch")
}

// SavedMealCreate creates a new saved meal for the authenticated user.
// Returns the newly created saved_meal_id on success.
//
// https://platform.fatsecret.com/docs/v1/saved_meal.create
func (c *FSProfileClient) SavedMealCreate(p SavedMealCreateParams) (*SavedMealCreateResponse, error) {
	params := map[string]string{
		"format":          "json",
		"saved_meal_name": p.SavedMealName,
	}
	if p.SavedMealDescription != "" {
		params["saved_meal_description"] = p.SavedMealDescription
	}
	if p.Meals != "" {
		params["meals"] = p.Meals
	}

	data, err := c.post(savedMealCreatePath, params)
	if err != nil {
		return nil, fmt.Errorf("SavedMealCreate: %w", err)
	}

	var resp SavedMealCreateResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("SavedMealCreate: decode response: %w", err)
	}

	return &resp, nil
}

