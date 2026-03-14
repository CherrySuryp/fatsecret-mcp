package fsclient

import "encoding/json"

const getRecentlyEatenPath string = "food/recently-eaten/v2"

type GetRecentlyEatenReq struct {
	Meal *MealType // optional
}

type GetRecentlyEatenResp struct {
	FoodID        string  `json:"food_id"`
	FoodName      string  `json:"food_name"`
	FoodType      string  `json:"food_type"`
	FoodUrl       string  `json:"food_url"`
	ServingID     string  `json:"serving_id"`
	NumberOfUnits float64 `json:"number_of_units,string"`
}

type getRecentlyEatenWrapper struct {
	Foods struct {
		Food []GetRecentlyEatenResp `json:"food"`
	} `json:"foods"`
}

func (c *Client) GetRecentlyEaten(payload GetRecentlyEatenReq) ([]GetRecentlyEatenResp, error) {
	params := map[string]string{}
	if payload.Meal != nil {
		params["meal"] = string(*payload.Meal)
	}

	resp, err := c.get(getRecentlyEatenPath, params)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var wrapper getRecentlyEatenWrapper
	err = json.Unmarshal(data, &wrapper)
	return wrapper.Foods.Food, err
}
