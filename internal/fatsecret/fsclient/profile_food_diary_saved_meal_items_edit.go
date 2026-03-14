package fsclient

import (
	"encoding/json"
	"fmt"
)

type SavedMealItemEditReq struct {
	SavedMealItemID   string
	SavedMealItemName *string
	NumberOfUnits     *float64
}

func (c *Client) SavedMealItemEdit(payload SavedMealItemEditReq) (SuccessResp, error) {
	var result SuccessResp

	params := map[string]string{
		"saved_meal_item_id": payload.SavedMealItemID,
	}
	if payload.SavedMealItemName != nil {
		params["saved_meal_item_name"] = *payload.SavedMealItemName
	}
	if payload.NumberOfUnits != nil {
		params["number_of_units"] = fmt.Sprintf("%.3f", *payload.NumberOfUnits)
	}

	resp, err := c.put(savedMealItemsPath, params)
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