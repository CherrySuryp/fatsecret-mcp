package fsprofileapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// https://platform.fatsecret.com/docs/v1/food_entry.create
const foodEntryCreatePath = "food-entries/v1"

// FoodEntryCreateParams holds the parameters for the food_entry.create endpoint.
type FoodEntryCreateParams struct {
	FoodID        string // required
	FoodEntryName string // required
	ServingID     string // required
	NumberOfUnits string // required; decimal as string
	Meal          string // required: "breakfast", "lunch", "dinner", or "other"
	Date          int    // required; days since 1970-01-01
}

// FoodEntryCreateResponse is the response from the food_entry.create endpoint.
type FoodEntryCreateResponse struct {
	FoodEntries FoodEntriesResult `json:"food_entries"`
}

// FoodEntryCreate adds a food item to the authenticated user's diary for the
// given date and meal. Returns the created entry with full nutritional data.
//
// https://platform.fatsecret.com/docs/v1/food_entry.create
func (c *FSProfileClient) FoodEntryCreate(p FoodEntryCreateParams) (*FoodEntryCreateResponse, error) {
	params := map[string]string{
		"format":          "json",
		"food_id":         p.FoodID,
		"food_entry_name": p.FoodEntryName,
		"serving_id":      p.ServingID,
		"number_of_units": p.NumberOfUnits,
		"meal":            p.Meal,
		"date":            strconv.Itoa(p.Date),
	}

	data, err := c.post(foodEntryCreatePath, params)
	if err != nil {
		return nil, fmt.Errorf("FoodEntryCreate: %w", err)
	}

	var resp FoodEntryCreateResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("FoodEntryCreate: decode response: %w", err)
	}

	return &resp, nil
}
