package fsprofileapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// https://platform.fatsecret.com/docs/v2/food_entries.get
const foodEntriesGetPath = "food-entries/v2"

// FoodEntriesGetParams holds the parameters for the food_entries.get endpoint.
type FoodEntriesGetParams struct {
	Date        int    // required; days since 1970-01-01
	FoodEntryID string // optional; filter to a specific entry
}

// FoodEntriesGetResponse is the response from the food_entries.get endpoint.
type FoodEntriesGetResponse struct {
	FoodEntries FoodEntriesResult `json:"food_entries"`
}

// FoodEntriesResult holds the list of food entries.
type FoodEntriesResult struct {
	Entries []FoodEntry `json:"food_entry"`
}

// FoodEntry is a single diary food entry with full nutritional breakdown.
type FoodEntry struct {
	FoodEntryID          string `json:"food_entry_id"`
	FoodEntryName        string `json:"food_entry_name"`
	FoodEntryDescription string `json:"food_entry_description"`
	DateInt              string `json:"date_int"`
	Meal                 string `json:"meal"`
	FoodID               string `json:"food_id"`
	ServingID            string `json:"serving_id"`
	NumberOfUnits        string `json:"number_of_units"`
	Calories             string `json:"calories"`
	Carbohydrate         string `json:"carbohydrate"`
	Protein              string `json:"protein"`
	Fat                  string `json:"fat"`
	SaturatedFat         string `json:"saturated_fat"`
	PolyunsaturatedFat   string `json:"polyunsaturated_fat"`
	MonounsaturatedFat   string `json:"monounsaturated_fat"`
	Cholesterol          string `json:"cholesterol"`
	Sodium               string `json:"sodium"`
	Potassium            string `json:"potassium"`
	Fiber                string `json:"fiber"`
	Sugar                string `json:"sugar"`
	VitaminA             string `json:"vitamin_a"`
	VitaminC             string `json:"vitamin_c"`
	Calcium              string `json:"calcium"`
	Iron                 string `json:"iron"`
}

// FoodEntriesGet returns food diary entries for the authenticated user on the
// given date, optionally filtered to a specific entry.
//
// https://platform.fatsecret.com/docs/v2/food_entries.get
func (c *FSProfileClient) FoodEntriesGet(p FoodEntriesGetParams) (*FoodEntriesGetResponse, error) {
	params := map[string]string{
		"format": "json",
		"date":   strconv.Itoa(p.Date),
	}
	if p.FoodEntryID != "" {
		params["food_entry_id"] = p.FoodEntryID
	}

	data, err := c.get(foodEntriesGetPath, params)
	if err != nil {
		return nil, fmt.Errorf("FoodEntriesGet: %w", err)
	}

	var resp FoodEntriesGetResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("FoodEntriesGet: decode response: %w", err)
	}

	return &resp, nil
}
