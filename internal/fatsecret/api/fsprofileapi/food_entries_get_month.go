package fsprofileapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// https://platform.fatsecret.com/docs/v2/food_entries.get_month
const foodEntriesGetMonthPath = "food-entries/month/v1"

// FoodEntriesGetMonthResponse is the response from the food_entries.get_month endpoint.
type FoodEntriesGetMonthResponse struct {
	Month FoodEntriesMonth `json:"month"`
}

// FoodEntriesMonth holds the monthly summary of food diary entries.
type FoodEntriesMonth struct {
	FromDateInt string         `json:"from_date_int"`
	ToDateInt   string         `json:"to_date_int"`
	Days        []FoodEntryDay `json:"day"`
}

// FoodEntryDay is a daily nutritional summary within a monthly report.
type FoodEntryDay struct {
	DateInt      string `json:"date_int"`
	Calories     string `json:"calories"`
	Carbohydrate string `json:"carbohydrate"`
	Protein      string `json:"protein"`
	Fat          string `json:"fat"`
}

// FoodEntriesGetMonth returns a monthly summary of food diary entries for the
// authenticated user. date is any day within the desired month (days since 1970-01-01).
//
// https://platform.fatsecret.com/docs/v2/food_entries.get_month
func (c *FSProfileClient) FoodEntriesGetMonth(date int) (*FoodEntriesGetMonthResponse, error) {
	params := map[string]string{
		"format": "json",
		"date":   strconv.Itoa(date),
	}

	data, err := c.get(foodEntriesGetMonthPath, params)
	if err != nil {
		return nil, fmt.Errorf("FoodEntriesGetMonth: %w", err)
	}

	var resp FoodEntriesGetMonthResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("FoodEntriesGetMonth: decode response: %w", err)
	}

	return &resp, nil
}
