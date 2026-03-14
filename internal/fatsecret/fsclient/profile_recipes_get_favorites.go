package fsclient

import "encoding/json"

const getRecipeFavoritesPath = "recipe/favorites/v2"

type GetRecipeFavoritesResp struct {
	RecipeID          string `json:"recipe_id"`
	RecipeName        string `json:"recipe_name"`
	RecipeUrl         string `json:"recipe_url"`
	RecipeDescription string `json:"recipe_description"`
	RecipeImage       string `json:"recipe_image"`
}

type getRecipeFavoritesWrapper struct {
	Recipes struct {
		Recipe []GetRecipeFavoritesResp `json:"recipe"`
	} `json:"recipes"`
}

func (c *Client) GetRecipeFavorites() ([]GetRecipeFavoritesResp, error) {
	resp, err := c.get(getRecipeFavoritesPath, make(map[string]string))
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var wrapper getRecipeFavoritesWrapper
	err = json.Unmarshal(data, &wrapper)
	return wrapper.Recipes.Recipe, err
}