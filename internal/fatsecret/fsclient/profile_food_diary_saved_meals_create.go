package fsclient

import (
	"encoding/json"
	"strings"
)

const savedMealsPath = "saved-meals/v1"

type SavedMealCreateReq struct {
	SavedMealName        string
	SavedMealDescription *string
	Meals                []MealType // optional; e.g. [MealTypeBreakfast, MealTypeLunch]
}

type SavedMealCreateResp struct {
	SavedMealID struct {
		Value string `json:"value"`
	} `json:"saved_meal_id"`
}

func (c *Client) SavedMealCreate(payload SavedMealCreateReq) (SavedMealCreateResp, error) {
	var result SavedMealCreateResp

	params := map[string]string{
		"saved_meal_name": payload.SavedMealName,
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

	resp, err := c.post(savedMealsPath, params)
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