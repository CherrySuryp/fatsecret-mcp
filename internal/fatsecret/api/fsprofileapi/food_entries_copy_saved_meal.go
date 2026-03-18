package fsprofileapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// https://platform.fatsecret.com/docs/v1/food_entries.copy_saved_meal
const foodEntriesCopySavedMealPath = "food-entries/copy/saved-meal/v1"

// FoodEntriesCopySavedMealParams holds the parameters for the food_entries.copy_saved_meal endpoint.
type FoodEntriesCopySavedMealParams struct {
	SavedMealID string // required
	Meal        string // required: "breakfast", "lunch", "dinner", or "other"
	Date        *int   // optional; days since 1970-01-01, defaults to current day
}

// FoodEntriesCopySavedMeal copies all items from a saved meal into the food
// diary for the authenticated user on the given date and meal slot.
//
// https://platform.fatsecret.com/docs/v1/food_entries.copy_saved_meal
func (c *FSProfileClient) FoodEntriesCopySavedMeal(p FoodEntriesCopySavedMealParams) (*SuccessResponse, error) {
	params := map[string]string{
		"format":        "json",
		"saved_meal_id": p.SavedMealID,
		"meal":          p.Meal,
	}
	if p.Date != nil {
		params["date"] = strconv.Itoa(*p.Date)
	}

	data, err := c.post(foodEntriesCopySavedMealPath, params)
	if err != nil {
		return nil, fmt.Errorf("FoodEntriesCopySavedMeal: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("FoodEntriesCopySavedMeal: decode response: %w", err)
	}

	return &resp, nil
}
