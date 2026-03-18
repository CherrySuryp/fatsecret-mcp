package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v2/recipes.get_favorites
const recipesGetFavoritesPath = "recipe/favorites/v2"

// RecipesGetFavoritesResponse is the response from the recipes.get_favorites endpoint.
type RecipesGetFavoritesResponse struct {
	Recipes RecipesGetFavoritesResult `json:"recipes"`
}

// RecipesGetFavoritesResult holds the list of favorite recipes.
type RecipesGetFavoritesResult struct {
	Recipes []FavoriteRecipeItem `json:"recipe"`
}

// FavoriteRecipeItem is a single favorite recipe entry.
type FavoriteRecipeItem struct {
	RecipeID          string `json:"recipe_id"`
	RecipeName        string `json:"recipe_name"`
	RecipeURL         string `json:"recipe_url"`
	RecipeDescription string `json:"recipe_description"`
	RecipeImage       string `json:"recipe_image"`
}

// RecipesGetFavorites returns the list of recipes saved as favorites for the
// authenticated user.
//
// https://platform.fatsecret.com/docs/v2/recipes.get_favorites
func (c *FSProfileClient) RecipesGetFavorites() (*RecipesGetFavoritesResponse, error) {
	params := map[string]string{"format": "json"}

	data, err := c.get(recipesGetFavoritesPath, params)
	if err != nil {
		return nil, fmt.Errorf("RecipesGetFavorites: %w", err)
	}

	var resp RecipesGetFavoritesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("RecipesGetFavorites: decode response: %w", err)
	}

	return &resp, nil
}

