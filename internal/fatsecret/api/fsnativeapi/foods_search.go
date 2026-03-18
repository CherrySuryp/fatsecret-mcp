package fsnativeapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const foodsSearchPath = "foods/search/v1"

// foodsSearchEnvelope is the private raw envelope used for JSON unmarshaling.
type foodsSearchEnvelope struct {
	Foods foodsSearchResultRaw `json:"foods"`
}

// FoodsSearchParams holds the parameters for the v1/foods.search endpoint.
// All fields are optional; SearchExpression is the primary filter.
type FoodsSearchParams struct {
	SearchExpression string
	PageNumber       *int
	MaxResults       *int
	// Premier-only: "weight" or "portion" — controls nutritional summary display.
	GenericDescription *string
	// Premier-only: filter results by region (e.g. "US").
	Region *string
	// Premier-only: return results in the specified language (requires Region).
	Language *string
}

// FoodsSearchResponse is the top-level response envelope for v1/foods.search.
type FoodsSearchResponse struct {
	Foods FoodsSearchResult `json:"foods"`
}

// FoodsSearchResult holds the pagination metadata and the food list.
// Numeric fields are strings in the FatSecret JSON response.
type FoodsSearchResult struct {
	MaxResults   string     `json:"max_results"`
	TotalResults string     `json:"total_results"`
	PageNumber   string     `json:"page_number"`
	Food         []FoodItem `json:"food"`
}

// FoodItem is a single food entry returned by the search.
type FoodItem struct {
	FoodID          string `json:"food_id"`
	FoodName        string `json:"food_name"`
	BrandName       string `json:"brand_name,omitempty"`
	FoodType        string `json:"food_type"`
	FoodURL         string `json:"food_url"`
	FoodDescription string `json:"food_description"`
}

// foodsSearchResultRaw is used internally to handle the single-vs-array quirk.
type foodsSearchResultRaw struct {
	MaxResults   string          `json:"max_results"`
	TotalResults string          `json:"total_results"`
	PageNumber   string          `json:"page_number"`
	Food         json.RawMessage `json:"food"`
}

// FoodsSearch calls the v1/foods.search endpoint and returns matching foods.
// Docs: https://platform.fatsecret.com/docs/v1/foods.search
func (c *FSNativeAPIClient) FoodsSearch(params FoodsSearchParams) (*FoodsSearchResponse, error) {
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
	if params.GenericDescription != nil {
		p["generic_description"] = *params.GenericDescription
	}
	if params.Region != nil {
		p["region"] = *params.Region
	}
	if params.Language != nil {
		p["language"] = *params.Language
	}

	body, err := c.get(foodsSearchPath, p)
	if err != nil {
		return nil, err
	}

	var raw foodsSearchEnvelope
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("fsnativeapi: unmarshal foods.search response: %w", err)
	}

	result := FoodsSearchResponse{
		Foods: FoodsSearchResult{
			MaxResults:   raw.Foods.MaxResults,
			TotalResults: raw.Foods.TotalResults,
			PageNumber:   raw.Foods.PageNumber,
		},
	}

	// FatSecret returns `food` as an object when there is exactly one result
	// and as an array otherwise.
	if len(raw.Foods.Food) > 0 {
		switch raw.Foods.Food[0] {
		case '[':
			if err := json.Unmarshal(raw.Foods.Food, &result.Foods.Food); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal food array: %w", err)
			}
		case '{':
			var single FoodItem
			if err := json.Unmarshal(raw.Foods.Food, &single); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal food object: %w", err)
			}
			result.Foods.Food = []FoodItem{single}
		}
	}

	return &result, nil
}
