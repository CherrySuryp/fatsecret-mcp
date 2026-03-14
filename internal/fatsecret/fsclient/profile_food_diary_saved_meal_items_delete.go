package fsclient

import "encoding/json"

type SavedMealItemDeleteReq struct {
	SavedMealItemID string
}

func (c *Client) SavedMealItemDelete(payload SavedMealItemDeleteReq) (SuccessResp, error) {
	var result SuccessResp

	params := map[string]string{
		"saved_meal_item_id": payload.SavedMealItemID,
	}

	resp, err := c.delete(savedMealItemsPath, params)
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