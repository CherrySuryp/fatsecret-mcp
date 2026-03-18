package fsnativeapi

import (
	"encoding/json"
	"fmt"
)

const naturalLanguageProcessingPath = "natural-language-processing/v1"

// NaturalLanguageProcessingParams holds the parameters for the
// natural-language-processing/v1 endpoint.
type NaturalLanguageProcessingParams struct {
	// UserInput is a description of what a user has eaten. Required. Max 1000 characters.
	UserInput string `json:"user_input"`
	// IncludeFoodData includes full food data in each response item when true.
	IncludeFoodData *bool `json:"include_food_data,omitempty"`
	// EatenFoods provides previously consumed foods to improve matching accuracy.
	EatenFoods []EatenFood `json:"eaten_foods,omitempty"`
	// Region filters results by region (e.g. "US", "FR"). Premier only.
	Region *string `json:"region,omitempty"`
	// Language specifies the response language. Requires Region. Premier only.
	Language *string `json:"language,omitempty"`
}

// NaturalLanguageProcessingResponse is the top-level response envelope for
// natural-language-processing/v1.
type NaturalLanguageProcessingResponse struct {
	FoodResponse []FoodResponseItem `json:"food_response"`
}

// NaturalLanguageProcessing calls the natural-language-processing/v1 endpoint,
// parsing a free-text food description into structured food and serving data.
// Docs: https://platform.fatsecret.com/docs/v1/natural.language.processing
func (c *FSNativeAPIClient) NaturalLanguageProcessing(params NaturalLanguageProcessingParams) (*NaturalLanguageProcessingResponse, error) {
	body, err := c.post(naturalLanguageProcessingPath, params)
	if err != nil {
		return nil, err
	}

	var result NaturalLanguageProcessingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("fsnativeapi: unmarshal natural.language.processing response: %w", err)
	}

	return &result, nil
}