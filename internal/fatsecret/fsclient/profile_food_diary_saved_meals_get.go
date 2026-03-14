package fsclient

import "encoding/json"

const savedMealsGetPath = "saved-meals/v2"

type SavedMealsGetReq struct {
	Meal *MealType // optional filter
}

type SavedMealsGetResp struct {
	SavedMealID          string `json:"saved_meal_id"`
	SavedMealName        string `json:"saved_meal_name"`
	SavedMealDescription string `json:"saved_meal_description"`
	Meals                string `json:"meals"`
}

type savedMealsGetWrapper struct {
	SavedMeals struct {
		SavedMeal []SavedMealsGetResp `json:"saved_meal"`
	} `json:"saved_meals"`
}

func (c *Client) SavedMealsGet(payload SavedMealsGetReq) ([]SavedMealsGetResp, error) {
	params := map[string]string{}
	if payload.Meal != nil {
		params["meal"] = string(*payload.Meal)
	}

	resp, err := c.get(savedMealsGetPath, params)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var wrapper savedMealsGetWrapper
	err = json.Unmarshal(data, &wrapper)
	return wrapper.SavedMeals.SavedMeal, err
}