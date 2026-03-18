package fsnativeapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const foodsSearchV5Path = "foods/search/v5"

// FoodsSearchV5Params holds the parameters for the foods/search/v5 endpoint.
type FoodsSearchV5Params struct {
	SearchExpression      string
	PageNumber            *int
	MaxResults            *int
	IncludeSubCategories  *bool
	IncludeFoodImages     *bool
	IncludeFoodAttributes *bool
	FlagDefaultServing    *bool
	// FoodType filters by "generic" or "brand". Defaults to no filter ("none").
	FoodType *string
	// Region filters results by region (e.g. "US", "FR"). Defaults to "US".
	Region *string
	// Language specifies the response language. Requires Region to be set.
	Language *string
}

// FoodsSearchV5Response is the top-level response envelope for foods/search/v5.
type FoodsSearchV5Response struct {
	FoodsSearch FoodsSearchV5Result `json:"foods_search"`
}

// FoodsSearchV5Result holds pagination metadata and the list of matched foods.
// Numeric fields are returned as strings by the FatSecret API.
type FoodsSearchV5Result struct {
	MaxResults   string          `json:"max_results"`
	TotalResults string          `json:"total_results"`
	PageNumber   string          `json:"page_number"`
	Results      []FoodV5Item    `json:"results"`
}

// FoodV5Item is a single food entry returned by the v5 search endpoint.
// It includes full serving data and optional images, attributes, and sub-categories.
type FoodV5Item struct {
	FoodID            string          `json:"food_id"`
	FoodName          string          `json:"food_name"`
	FoodType          string          `json:"food_type"`
	FoodURL           string          `json:"food_url"`
	BrandName         string          `json:"brand_name,omitempty"`
	FoodSubCategories *SubCategoryList `json:"food_sub_categories,omitempty"`
	FoodImages        *FoodImageList   `json:"food_images,omitempty"`
	FoodAttributes    *FoodAttributes  `json:"food_attributes,omitempty"`
	Servings          struct {
		Serving []Serving `json:"serving"`
	} `json:"servings"`
}

type foodsSearchV5Envelope struct {
	FoodsSearch foodsSearchV5ResultRaw `json:"foods_search"`
}

type foodsSearchV5ResultRaw struct {
	MaxResults string          `json:"max_results"`
	TotalResults string        `json:"total_results"`
	PageNumber string          `json:"page_number"`
	Results    json.RawMessage `json:"results"`
}

// FoodsSearchV5 calls the foods/search/v5 endpoint and returns matching foods
// with full nutritional detail per serving.
// Docs: https://platform.fatsecret.com/docs/v5/foods.search
func (c *FSNativeAPIClient) FoodsSearchV5(params FoodsSearchV5Params) (*FoodsSearchV5Response, error) {
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
	if params.IncludeSubCategories != nil {
		p["include_sub_categories"] = strconv.FormatBool(*params.IncludeSubCategories)
	}
	if params.IncludeFoodImages != nil {
		p["include_food_images"] = strconv.FormatBool(*params.IncludeFoodImages)
	}
	if params.IncludeFoodAttributes != nil {
		p["include_food_attributes"] = strconv.FormatBool(*params.IncludeFoodAttributes)
	}
	if params.FlagDefaultServing != nil {
		p["flag_default_serving"] = strconv.FormatBool(*params.FlagDefaultServing)
	}
	if params.FoodType != nil {
		p["food_type"] = *params.FoodType
	}
	if params.Region != nil {
		p["region"] = *params.Region
	}
	if params.Language != nil {
		p["language"] = *params.Language
	}

	body, err := c.get(foodsSearchV5Path, p)
	if err != nil {
		return nil, err
	}

	var raw foodsSearchV5Envelope
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("fsnativeapi: unmarshal foods.search v5 response: %w", err)
	}

	result := FoodsSearchV5Response{
		FoodsSearch: FoodsSearchV5Result{
			MaxResults:   raw.FoodsSearch.MaxResults,
			TotalResults: raw.FoodsSearch.TotalResults,
			PageNumber:   raw.FoodsSearch.PageNumber,
		},
	}

	// `results` may be an object (single food) or array (multiple foods).
	if len(raw.FoodsSearch.Results) > 0 {
		switch raw.FoodsSearch.Results[0] {
		case '[':
			if err := json.Unmarshal(raw.FoodsSearch.Results, &result.FoodsSearch.Results); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal foods v5 array: %w", err)
			}
		case '{':
			var single FoodV5Item
			if err := json.Unmarshal(raw.FoodsSearch.Results, &single); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal foods v5 object: %w", err)
			}
			result.FoodsSearch.Results = []FoodV5Item{single}
		}
	}

	return &result, nil
}