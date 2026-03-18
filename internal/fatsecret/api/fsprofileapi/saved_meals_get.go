package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v2/saved_meals.get
const savedMealsGetPath = "saved-meals/v2"

// SavedMealsGetResponse is the response from the saved_meals.get endpoint.
type SavedMealsGetResponse struct {
	SavedMeals SavedMealsResult `json:"saved_meals"`
}

// SavedMealsResult holds the list of saved meals.
type SavedMealsResult struct {
	Meals []SavedMeal `json:"saved_meal"`
}

// SavedMeal is a single saved meal entry.
type SavedMeal struct {
	SavedMealID          string `json:"saved_meal_id"`
	SavedMealName        string `json:"saved_meal_name"`
	SavedMealDescription string `json:"saved_meal_description"`
	Meals                string `json:"meals"` // comma-separated meal types
}

// SavedMealsGet returns the saved meals for the authenticated user,
// optionally filtered by meal type.
//
// https://platform.fatsecret.com/docs/v2/saved_meals.get
func (c *FSProfileClient) SavedMealsGet(meal string) (*SavedMealsGetResponse, error) {
	params := map[string]string{"format": "json"}
	if meal != "" {
		params["meal"] = meal
	}

	data, err := c.get(savedMealsGetPath, params)
	if err != nil {
		return nil, fmt.Errorf("SavedMealsGet: %w", err)
	}

	var resp SavedMealsGetResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("SavedMealsGet: decode response: %w", err)
	}

	return &resp, nil
}

