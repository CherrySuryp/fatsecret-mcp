package fsprofileapi

import "encoding/json"

const savedMealItemsGetPath = "saved-meals/item/v2"

type SavedMealItemsGetReq struct {
	SavedMealID string
}

type SavedMealItemsGetResp struct {
	SavedMealItemID   string  `json:"saved_meal_item_id"`
	FoodID            string  `json:"food_id"`
	SavedMealItemName string  `json:"saved_meal_item_name"`
	ServingID         string  `json:"serving_id"`
	NumberOfUnits     float64 `json:"number_of_units,string"`
}

type savedMealItemsGetWrapper struct {
	SavedMealItems struct {
		SavedMealItem []SavedMealItemsGetResp `json:"saved_meal_item"`
	} `json:"saved_meal_items"`
}

func (c *Client) SavedMealItemsGet(payload SavedMealItemsGetReq) ([]SavedMealItemsGetResp, error) {
	params := map[string]string{
		"saved_meal_id": payload.SavedMealID,
	}

	resp, err := c.get(savedMealItemsGetPath, params)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var wrapper savedMealItemsGetWrapper
	err = json.Unmarshal(data, &wrapper)
	return wrapper.SavedMealItems.SavedMealItem, err
}
