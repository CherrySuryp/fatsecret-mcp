package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v1/recipe.delete_favorite
const recipeDeleteFavoritePath = "recipe/favorites/v1"

// RecipeDeleteFavorite removes a recipe from the authenticated user's favorites.
//
// https://platform.fatsecret.com/docs/v1/recipe.delete_favorite
func (c *FSProfileClient) RecipeDeleteFavorite(recipeID string) (*SuccessResponse, error) {
	params := map[string]string{
		"format":    "json",
		"recipe_id": recipeID,
	}

	data, err := c.delete(recipeDeleteFavoritePath, params)
	if err != nil {
		return nil, fmt.Errorf("RecipeDeleteFavorite: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("RecipeDeleteFavorite: decode response: %w", err)
	}

	return &resp, nil
}

