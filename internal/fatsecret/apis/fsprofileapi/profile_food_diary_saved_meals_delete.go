package fsprofileapi

import "encoding/json"

type SavedMealDeleteReq struct {
	SavedMealID string
}

func (c *Client) SavedMealDelete(payload SavedMealDeleteReq) (SuccessResp, error) {
	var result SuccessResp

	params := map[string]string{
		"saved_meal_id": payload.SavedMealID,
	}

	resp, err := c.delete(savedMealsPath, params)
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
