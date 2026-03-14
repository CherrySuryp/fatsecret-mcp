package fsclient

import (
	"encoding/json"
	"fmt"
)

const savedMealItemsPath = "saved-meals/item/v1"

type SavedMealItemAddReq struct {
	SavedMealID       string
	FoodID            string
	SavedMealItemName string
	ServingID         string
	NumberOfUnits     float64
}

type SavedMealItemAddResp struct {
	SavedMealItemID struct {
		Value string `json:"value"`
	} `json:"saved_meal_item_id"`
}

func (c *Client) SavedMealItemAdd(payload SavedMealItemAddReq) (SavedMealItemAddResp, error) {
	var result SavedMealItemAddResp

	params := map[string]string{
		"saved_meal_id":        payload.SavedMealID,
		"food_id":              payload.FoodID,
		"saved_meal_item_name": payload.SavedMealItemName,
		"serving_id":           payload.ServingID,
		"number_of_units":      fmt.Sprintf("%.3f", payload.NumberOfUnits),
	}

	resp, err := c.post(savedMealItemsPath, params)
	if err != nil {
		return result, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(data, &result)
	return result, err
}