package fsclient

import "encoding/json"

const addRecipeFavoritePath = "recipe/favorites/v1"

type RecipeAddFavoriteReq struct {
	RecipeID string
}

func (c *Client) RecipeAddFavorite(payload RecipeAddFavoriteReq) (SuccessResp, error) {
	var result SuccessResp

	params := map[string]string{
		"recipe_id": payload.RecipeID,
	}

	resp, err := c.post(addRecipeFavoritePath, params)
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