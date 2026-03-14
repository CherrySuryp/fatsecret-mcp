package fsclient

import "encoding/json"

const getAllFavoritesPath string = "food/favorites/v2"

type GetAllFavoritesResp struct {
	FoodID          string  `json:"food_id"`
	FoodName        string  `json:"food_name"`
	FoodType        string  `json:"food_type"`
	FoodUrl         string  `json:"food_url"`
	FoodDescription string  `json:"food_description"`
	ServingID       string  `json:"serving_id"`
	NumberOfUnits   float64 `json:"number_of_units,string"`
}

type getAllFavoritesWrapper struct {
	Foods struct {
		Food []GetAllFavoritesResp `json:"food"`
	} `json:"foods"`
}

func (c *Client) GetAllFavorites() ([]GetAllFavoritesResp, error) {
	resp, err := c.get(getAllFavoritesPath, make(map[string]string))
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var wrapper getAllFavoritesWrapper
	err = json.Unmarshal(data, &wrapper)
	return wrapper.Foods.Food, err
}
