package fsprofileapi

import (
	"encoding/json"
	"strings"
)

type SavedMealEditReq struct {
	SavedMealID          string
	SavedMealName        *string
	SavedMealDescription *string
	Meals                []MealType // optional
}

func (c *Client) SavedMealEdit(payload SavedMealEditReq) (SuccessResp, error) {
	var result SuccessResp

	params := map[string]string{
		"saved_meal_id": payload.SavedMealID,
	}
	if payload.SavedMealName != nil {
		params["saved_meal_name"] = *payload.SavedMealName
	}
	if payload.SavedMealDescription != nil {
		params["saved_meal_description"] = *payload.SavedMealDescription
	}
	if len(payload.Meals) > 0 {
		meals := make([]string, len(payload.Meals))
		for i, m := range payload.Meals {
			meals[i] = string(m)
		}
		params["meals"] = strings.Join(meals, ",")
	}

	resp, err := c.put(savedMealsPath, params)
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
