package fsclient

import "encoding/json"

const deleteRecipeFavoritePath = "recipe/favorites/v1"

type RecipeDeleteFavoriteReq struct {
	RecipeID string
}

func (c *Client) RecipeDeleteFavorite(payload RecipeDeleteFavoriteReq) (SuccessResp, error) {
	var result SuccessResp

	params := map[string]string{
		"recipe_id": payload.RecipeID,
	}

	resp, err := c.delete(deleteRecipeFavoritePath, params)
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