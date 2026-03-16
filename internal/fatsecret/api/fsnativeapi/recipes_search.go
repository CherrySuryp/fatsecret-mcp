package fsnativeapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const recipesSearchPath = "recipes/search/v3"

// RecipesSearchParams holds the parameters for the recipes/search/v3 endpoint.
type RecipesSearchParams struct {
	SearchExpression string
	PageNumber       *int
	MaxResults       *int
	// RecipeTypes is a comma-separated list of recipe type names to filter by.
	RecipeTypes *string
	// RecipeTypesMatchAll requires all listed types to match when true (default: any).
	RecipeTypesMatchAll *bool
	// MustHaveImages limits results to recipes with images.
	MustHaveImages *bool
	// CaloriesFrom and CaloriesTo filter by calorie range (kcal per serving).
	CaloriesFrom *int
	CaloriesTo   *int
	// CarbPercentageFrom/To filter by carbohydrate percentage of total calories.
	CarbPercentageFrom *int
	CarbPercentageTo   *int
	// ProteinPercentageFrom/To filter by protein percentage of total calories.
	ProteinPercentageFrom *int
	ProteinPercentageTo   *int
	// FatPercentageFrom/To filter by fat percentage of total calories.
	FatPercentageFrom *int
	FatPercentageTo   *int
	// PrepTimeFrom/To filter by preparation time range in minutes.
	PrepTimeFrom *int
	PrepTimeTo   *int
	// SortBy controls result ordering. Valid values: "newest", "oldest",
	// "caloriesPerServingAscending", "caloriesPerServingDescending".
	SortBy *string
	// Region filters results by region (e.g. "US", "FR"). Defaults to "US". Premier only.
	Region *string
}

// RecipesSearchResponse is the top-level response envelope for recipes/search/v3.
type RecipesSearchResponse struct {
	Recipes RecipesSearchResult `json:"recipes"`
}

// RecipesSearchResult holds pagination metadata and the list of matched recipes.
type RecipesSearchResult struct {
	MaxResults   string            `json:"max_results"`
	TotalResults string            `json:"total_results"`
	PageNumber   string            `json:"page_number"`
	Recipe       []RecipeSearchItem `json:"recipe"`
}

// RecipeSearchItem is a summary recipe entry returned by recipes/search/v3.
type RecipeSearchItem struct {
	RecipeID          string                   `json:"recipe_id"`
	RecipeName        string                   `json:"recipe_name"`
	RecipeDescription string                   `json:"recipe_description,omitempty"`
	RecipeImage       string                   `json:"recipe_image,omitempty"`
	RecipeNutrition   RecipeSearchNutrition    `json:"recipe_nutrition"`
	RecipeIngredients *RecipeSearchIngredients `json:"recipe_ingredients,omitempty"`
	RecipeTypes       *RecipeSearchTypes       `json:"recipe_types,omitempty"`
}

// RecipeSearchNutrition holds per-serving macronutrient totals for a search result recipe.
type RecipeSearchNutrition struct {
	Calories     string `json:"calories"`
	Carbohydrate string `json:"carbohydrate"`
	Protein      string `json:"protein"`
	Fat          string `json:"fat"`
}

// RecipeSearchIngredients wraps the ingredient string array in a search result.
type RecipeSearchIngredients struct {
	Ingredient []string `json:"ingredient"`
}

// RecipeSearchTypes wraps the recipe type string array in a search result.
type RecipeSearchTypes struct {
	RecipeType []string `json:"recipe_type"`
}

type recipesSearchEnvelope struct {
	Recipes recipesSearchResultRaw `json:"recipes"`
}

type recipesSearchResultRaw struct {
	MaxResults   string          `json:"max_results"`
	TotalResults string          `json:"total_results"`
	PageNumber   string          `json:"page_number"`
	Recipe       json.RawMessage `json:"recipe"`
}

// RecipesSearch calls the recipes/search/v3 endpoint and returns matching recipes.
// Docs: https://platform.fatsecret.com/docs/v3/recipes.search
func (c *FSNativeAPIClient) RecipesSearch(params RecipesSearchParams) (*RecipesSearchResponse, error) {
	p := map[string]string{
		"format": "json",
	}
	if params.SearchExpression != "" {
		p["search_expression"] = params.SearchExpression
	}
	if params.PageNumber != nil {
		p["page_number"] = strconv.Itoa(*params.PageNumber)
	}
	if params.MaxResults != nil {
		p["max_results"] = strconv.Itoa(*params.MaxResults)
	}
	if params.RecipeTypes != nil {
		p["recipe_types"] = *params.RecipeTypes
	}
	if params.RecipeTypesMatchAll != nil {
		p["recipe_types_matchall"] = strconv.FormatBool(*params.RecipeTypesMatchAll)
	}
	if params.MustHaveImages != nil {
		p["must_have_images"] = strconv.FormatBool(*params.MustHaveImages)
	}
	if params.CaloriesFrom != nil {
		p["calories.from"] = strconv.Itoa(*params.CaloriesFrom)
	}
	if params.CaloriesTo != nil {
		p["calories.to"] = strconv.Itoa(*params.CaloriesTo)
	}
	if params.CarbPercentageFrom != nil {
		p["carb_percentage.from"] = strconv.Itoa(*params.CarbPercentageFrom)
	}
	if params.CarbPercentageTo != nil {
		p["carb_percentage.to"] = strconv.Itoa(*params.CarbPercentageTo)
	}
	if params.ProteinPercentageFrom != nil {
		p["protein_percentage.from"] = strconv.Itoa(*params.ProteinPercentageFrom)
	}
	if params.ProteinPercentageTo != nil {
		p["protein_percentage.to"] = strconv.Itoa(*params.ProteinPercentageTo)
	}
	if params.FatPercentageFrom != nil {
		p["fat_percentage.from"] = strconv.Itoa(*params.FatPercentageFrom)
	}
	if params.FatPercentageTo != nil {
		p["fat_percentage.to"] = strconv.Itoa(*params.FatPercentageTo)
	}
	if params.PrepTimeFrom != nil {
		p["prep_time.from"] = strconv.Itoa(*params.PrepTimeFrom)
	}
	if params.PrepTimeTo != nil {
		p["prep_time.to"] = strconv.Itoa(*params.PrepTimeTo)
	}
	if params.SortBy != nil {
		p["sort_by"] = *params.SortBy
	}
	if params.Region != nil {
		p["region"] = *params.Region
	}

	body, err := c.get(recipesSearchPath, p)
	if err != nil {
		return nil, err
	}

	var raw recipesSearchEnvelope
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("fsnativeapi: unmarshal recipes.search response: %w", err)
	}

	result := RecipesSearchResponse{
		Recipes: RecipesSearchResult{
			MaxResults:   raw.Recipes.MaxResults,
			TotalResults: raw.Recipes.TotalResults,
			PageNumber:   raw.Recipes.PageNumber,
		},
	}

	// FatSecret returns `recipe` as an object when there is exactly one result
	// and as an array otherwise.
	if len(raw.Recipes.Recipe) > 0 {
		switch raw.Recipes.Recipe[0] {
		case '[':
			if err := json.Unmarshal(raw.Recipes.Recipe, &result.Recipes.Recipe); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal recipe array: %w", err)
			}
		case '{':
			var single RecipeSearchItem
			if err := json.Unmarshal(raw.Recipes.Recipe, &single); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal recipe object: %w", err)
			}
			result.Recipes.Recipe = []RecipeSearchItem{single}
		}
	}

	return &result, nil
}