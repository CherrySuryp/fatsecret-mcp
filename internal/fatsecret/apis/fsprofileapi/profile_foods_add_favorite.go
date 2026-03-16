package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

const addFavoritePath = "food/favorite/v1"

type FoodAddFavoriteReq struct {
	FoodID        string
	ServingID     *string
	NumberOfUnits *uint
}

func (c *Client) FoodAddFavorite(payload FoodAddFavoriteReq) (SuccessResp, error) {
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

	resp, err := c.post(addFavoritePath, params)
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
