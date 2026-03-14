package fsclient

import (
	"encoding/json"
	"fmt"
)

const deleteFavoritePath = "food/favorite/v1"

type FoodDeleteFavoriteReq struct {
	FoodID        string
	ServingID     *string
	NumberOfUnits *uint
}

func (c *Client) FoodDeleteFavorite(payload FoodDeleteFavoriteReq) (SuccessResp, error) {
	var result SuccessResp

	params := map[string]string{
		"food_id": payload.FoodID,
	}
	if payload.ServingID != nil {
		params["serving_id"] = *payload.ServingID
	}
	if payload.NumberOfUnits != nil {
		params["number_of_units"] = fmt.Sprintf("%.3f", float64(*payload.NumberOfUnits))
	}

	resp, err := c.delete(deleteFavoritePath, params)
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
