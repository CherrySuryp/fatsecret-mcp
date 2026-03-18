package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v1/recipe.add_favorite
const recipeAddFavoritePath = "recipe/favorites/v1"

// RecipeAddFavorite adds a recipe to the authenticated user's favorites.
//
// https://platform.fatsecret.com/docs/v1/recipe.add_favorite
func (c *FSProfileClient) RecipeAddFavorite(recipeID string) (*SuccessResponse, error) {
	params := map[string]string{
		"format":    "json",
		"recipe_id": recipeID,
	}

	data, err := c.post(recipeAddFavoritePath, params)
	if err != nil {
		return nil, fmt.Errorf("RecipeAddFavorite: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("RecipeAddFavorite: decode response: %w", err)
	}

	return &resp, nil
}

