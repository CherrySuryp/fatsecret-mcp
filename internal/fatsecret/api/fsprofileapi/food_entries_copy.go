package fsprofileapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// https://platform.fatsecret.com/docs/v1/food_entries.copy
const foodEntriesCopyPath = "food-entries/copy/v1"

// FoodEntriesCopyParams holds the parameters for the food_entries.copy endpoint.
type FoodEntriesCopyParams struct {
	FromDate int    // required; days since 1970-01-01
	ToDate   *int   // optional; defaults to current day
	Meal     string // optional: "breakfast", "lunch", "dinner", or "other"
}

// FoodEntriesCopy copies food entries from one date to another for the
// authenticated user.
//
// https://platform.fatsecret.com/docs/v1/food_entries.copy
func (c *FSProfileClient) FoodEntriesCopy(p FoodEntriesCopyParams) (*SuccessResponse, error) {
	params := map[string]string{
		"format":    "json",
		"from_date": strconv.Itoa(p.FromDate),
	}
	if p.ToDate != nil {
		params["to_date"] = strconv.Itoa(*p.ToDate)
	}
	if p.Meal != "" {
		params["meal"] = p.Meal
	}

	data, err := c.post(foodEntriesCopyPath, params)
	if err != nil {
		return nil, fmt.Errorf("FoodEntriesCopy: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("FoodEntriesCopy: decode response: %w", err)
	}

	return &resp, nil
}
