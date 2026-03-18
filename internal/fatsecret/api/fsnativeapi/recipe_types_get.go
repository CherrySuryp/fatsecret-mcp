package fsnativeapi

import (
	"encoding/json"
	"fmt"
)

const recipeTypesGetPath = "recipe-types/v2"

// RecipeTypesGetParams holds the parameters for the recipe-types/v2 endpoint.
type RecipeTypesGetParams struct {
	// Region filters results by region (e.g. "US", "FR"). Defaults to "US". Premier only.
	Region *string
	// Language specifies the response language. Requires Region to be set. Premier only.
	Language *string
}

// RecipeTypesGetResponse is the top-level response envelope for recipe-types/v2.
type RecipeTypesGetResponse struct {
	RecipeTypes RecipeTypesResult `json:"recipe_types"`
}

// RecipeTypesResult holds the list of recipe type name strings.
type RecipeTypesResult struct {
	RecipeType []string `json:"recipe_type"`
}

type recipeTypesEnvelope struct {
	RecipeTypes recipeTypesResultRaw `json:"recipe_types"`
}

type recipeTypesResultRaw struct {
	RecipeType json.RawMessage `json:"recipe_type"`
}

// RecipeTypesGet calls the recipe-types/v2 endpoint and returns all recipe type names.
// Docs: https://platform.fatsecret.com/docs/v2/recipe_types.get
func (c *FSNativeAPIClient) RecipeTypesGet(params RecipeTypesGetParams) (*RecipeTypesGetResponse, error) {
	p := map[string]string{
		"format": "json",
	}
	if params.Region != nil {
		p["region"] = *params.Region
	}
	if params.Language != nil {
		p["language"] = *params.Language
	}

	body, err := c.get(recipeTypesGetPath, p)
	if err != nil {
		return nil, err
	}

	var raw recipeTypesEnvelope
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("fsnativeapi: unmarshal recipe_types.get response: %w", err)
	}

	result := RecipeTypesGetResponse{}

	if len(raw.RecipeTypes.RecipeType) > 0 {
		switch raw.RecipeTypes.RecipeType[0] {
		case '[':
			if err := json.Unmarshal(raw.RecipeTypes.RecipeType, &result.RecipeTypes.RecipeType); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal recipe_type array: %w", err)
			}
		case '"':
			var single string
			if err := json.Unmarshal(raw.RecipeTypes.RecipeType, &single); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal recipe_type string: %w", err)
			}
			result.RecipeTypes.RecipeType = []string{single}
		}
	}

	return &result, nil
}
