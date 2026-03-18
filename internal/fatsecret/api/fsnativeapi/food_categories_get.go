package fsnativeapi

import (
	"encoding/json"
	"fmt"
)

const foodCategoriesGetPath = "food-categories/v2"

// FoodCategoriesGetParams holds the parameters for the food-categories/v2 endpoint.
type FoodCategoriesGetParams struct {
	// Region filters results by region (e.g. "US", "FR"). Defaults to "US".
	Region *string
	// Language specifies the response language. Requires Region to be set.
	Language *string
}

// FoodCategoriesGetResponse is the top-level response envelope for food-categories/v2.
type FoodCategoriesGetResponse struct {
	FoodCategories FoodCategoriesResult `json:"food_categories"`
}

// FoodCategoriesResult holds the list of food categories.
type FoodCategoriesResult struct {
	FoodCategory []FoodCategory `json:"food_category"`
}

// FoodCategory is a single food category entry.
type FoodCategory struct {
	FoodCategoryID          string `json:"food_category_id"`
	FoodCategoryName        string `json:"food_category_name"`
	FoodCategoryDescription string `json:"food_category_description"`
}

type foodCategoriesEnvelope struct {
	FoodCategories foodCategoriesResultRaw `json:"food_categories"`
}

type foodCategoriesResultRaw struct {
	FoodCategory json.RawMessage `json:"food_category"`
}

// FoodCategoriesGet calls the food-categories/v2 endpoint and returns all food categories.
// Docs: https://platform.fatsecret.com/docs/v2/food_categories.get
func (c *FSNativeAPIClient) FoodCategoriesGet(params FoodCategoriesGetParams) (*FoodCategoriesGetResponse, error) {
	p := map[string]string{
		"format": "json",
	}
	if params.Region != nil {
		p["region"] = *params.Region
	}
	if params.Language != nil {
		p["language"] = *params.Language
	}

	body, err := c.get(foodCategoriesGetPath, p)
	if err != nil {
		return nil, err
	}

	var raw foodCategoriesEnvelope
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("fsnativeapi: unmarshal food_categories.get response: %w", err)
	}

	result := FoodCategoriesGetResponse{}

	// FatSecret returns `food_category` as an object when there is exactly one result
	// and as an array otherwise.
	if len(raw.FoodCategories.FoodCategory) > 0 {
		switch raw.FoodCategories.FoodCategory[0] {
		case '[':
			if err := json.Unmarshal(raw.FoodCategories.FoodCategory, &result.FoodCategories.FoodCategory); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal food_category array: %w", err)
			}
		case '{':
			var single FoodCategory
			if err := json.Unmarshal(raw.FoodCategories.FoodCategory, &single); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal food_category object: %w", err)
			}
			result.FoodCategories.FoodCategory = []FoodCategory{single}
		}
	}

	return &result, nil
}