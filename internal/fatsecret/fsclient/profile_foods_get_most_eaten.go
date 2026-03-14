package fsclient

import "encoding/json"

const getMostEatenPath string = "food/most-eaten/v2"

type GetMostEatenReq struct {
	Meal *MealType // optional
}

type GetMostEatenResp struct {
	FoodID        string  `json:"food_id"`
	FoodName      string  `json:"food_name"`
	FoodType      string  `json:"food_type"`
	FoodUrl       string  `json:"food_url"`
	ServingID     string  `json:"serving_id"`
	NumberOfUnits float64 `json:"number_of_units,string"`
}

type getMostEatenWrapper struct {
	Foods struct {
		Food []GetMostEatenResp `json:"food"`
	} `json:"foods"`
}

func (c *Client) GetMostEaten(payload GetMostEatenReq) ([]GetMostEatenResp, error) {
	params := map[string]string{}
	if payload.Meal != nil {
		params["meal"] = string(*payload.Meal)
	}

	resp, err := c.get(getMostEatenPath, params)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var wrapper getMostEatenWrapper
	err = json.Unmarshal(data, &wrapper)
	return wrapper.Foods.Food, err
}
