package fsnativeapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const foodSubCategoriesGetPath = "food-sub-categories/v2"

// FoodSubCategoriesGetParams holds the parameters for the food-sub-categories/v2 endpoint.
type FoodSubCategoriesGetParams struct {
	// FoodCategoryID is the unique identifier of the parent food category. Required.
	FoodCategoryID int64
	// Region filters results by region (e.g. "US", "FR"). Defaults to "US".
	Region *string
	// Language specifies the response language. Requires Region to be set.
	Language *string
}

// FoodSubCategoriesGetResponse is the top-level response envelope for food-sub-categories/v2.
type FoodSubCategoriesGetResponse struct {
	FoodSubCategories FoodSubCategoriesResult `json:"food_sub_categories"`
}

// FoodSubCategoriesResult holds the list of sub-category name strings.
type FoodSubCategoriesResult struct {
	FoodSubCategory []string `json:"food_sub_category"`
}

type foodSubCategoriesEnvelope struct {
	FoodSubCategories foodSubCategoriesResultRaw `json:"food_sub_categories"`
}

type foodSubCategoriesResultRaw struct {
	FoodSubCategory json.RawMessage `json:"food_sub_category"`
}

// FoodSubCategoriesGet calls the food-sub-categories/v2 endpoint and returns all
// sub-categories for the given food category.
// Docs: https://platform.fatsecret.com/docs/v2/food_sub_categories.get
func (c *FSNativeAPIClient) FoodSubCategoriesGet(params FoodSubCategoriesGetParams) (*FoodSubCategoriesGetResponse, error) {
	p := map[string]string{
		"format":           "json",
		"food_category_id": strconv.FormatInt(params.FoodCategoryID, 10),
	}
	if params.Region != nil {
		p["region"] = *params.Region
	}
	if params.Language != nil {
		p["language"] = *params.Language
	}

	body, err := c.get(foodSubCategoriesGetPath, p)
	if err != nil {
		return nil, err
	}

	var raw foodSubCategoriesEnvelope
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("fsnativeapi: unmarshal food_sub_categories.get response: %w", err)
	}

	result := FoodSubCategoriesGetResponse{}

	// FatSecret returns `food_sub_category` as a string when there is exactly one result
	// and as an array otherwise.
	if len(raw.FoodSubCategories.FoodSubCategory) > 0 {
		switch raw.FoodSubCategories.FoodSubCategory[0] {
		case '[':
			if err := json.Unmarshal(raw.FoodSubCategories.FoodSubCategory, &result.FoodSubCategories.FoodSubCategory); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal food_sub_category array: %w", err)
			}
		case '"':
			var single string
			if err := json.Unmarshal(raw.FoodSubCategories.FoodSubCategory, &single); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal food_sub_category string: %w", err)
			}
			result.FoodSubCategories.FoodSubCategory = []string{single}
		}
	}

	return &result, nil
}