package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v1/saved_meal_item.add
const savedMealItemAddPath = "saved-meals/item/v1"

// SavedMealItemAddResponse is the response from the saved_meal_item.add endpoint.
type SavedMealItemAddResponse struct {
	SavedMealItemID struct {
		Value string `json:"value"`
	} `json:"saved_meal_item_id"`
}

// SavedMealItemAdd adds a food item to an existing saved meal.
// Returns the newly created saved_meal_item_id on success.
//
// https://platform.fatsecret.com/docs/v1/saved_meal_item.add
func (c *FSProfileClient) SavedMealItemAdd(savedMealID, foodID, savedMealItemName, servingID, numberOfUnits string) (*SavedMealItemAddResponse, error) {
	params := map[string]string{
		"format":               "json",
		"saved_meal_id":        savedMealID,
		"food_id":              foodID,
		"saved_meal_item_name": savedMealItemName,
		"serving_id":           servingID,
		"number_of_units":      numberOfUnits,
	}

	data, err := c.post(savedMealItemAddPath, params)
	if err != nil {
		return nil, fmt.Errorf("SavedMealItemAdd: %w", err)
	}

	var resp SavedMealItemAddResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("SavedMealItemAdd: decode response: %w", err)
	}

	return &resp, nil
}

