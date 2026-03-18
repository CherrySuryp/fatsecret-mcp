package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v2/saved_meal_items.get
const savedMealItemsGetPath = "saved-meals/item/v2"

// SavedMealItemsGetResponse is the response from the saved_meal_items.get endpoint.
type SavedMealItemsGetResponse struct {
	SavedMealItems SavedMealItemsResult `json:"saved_meal_items"`
}

// SavedMealItemsResult holds the list of saved meal items.
type SavedMealItemsResult struct {
	Items []SavedMealItem `json:"saved_meal_item"`
}

// SavedMealItem is a single item within a saved meal.
type SavedMealItem struct {
	SavedMealItemID   string `json:"saved_meal_item_id"`
	SavedMealItemName string `json:"saved_meal_item_name"`
	FoodID            string `json:"food_id"`
	ServingID         string `json:"serving_id"`
	NumberOfUnits     string `json:"number_of_units"`
}

// SavedMealItemsGet returns all items belonging to the given saved meal.
//
// https://platform.fatsecret.com/docs/v2/saved_meal_items.get
func (c *FSProfileClient) SavedMealItemsGet(savedMealID string) (*SavedMealItemsGetResponse, error) {
	params := map[string]string{
		"format":        "json",
		"saved_meal_id": savedMealID,
	}

	data, err := c.get(savedMealItemsGetPath, params)
	if err != nil {
		return nil, fmt.Errorf("SavedMealItemsGet: %w", err)
	}

	var resp SavedMealItemsGetResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("SavedMealItemsGet: decode response: %w", err)
	}

	return &resp, nil
}

