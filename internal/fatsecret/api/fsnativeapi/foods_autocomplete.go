package fsnativeapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const foodsAutocompletePath = "food/autocomplete/v2"

// FoodsAutocompleteParams holds the parameters for the food/autocomplete/v2 endpoint.
type FoodsAutocompleteParams struct {
	// Expression is the search term to autocomplete. Required.
	Expression string
	// MaxResults limits the number of suggestions returned. Default 4, max 10.
	MaxResults *int
	// Region filters results by region (e.g. "US", "FR"). Defaults to "US".
	Region *string
}

// FoodsAutocompleteResponse is the top-level response envelope for food/autocomplete/v2.
type FoodsAutocompleteResponse struct {
	Suggestions SuggestionsResult `json:"suggestions"`
}

// SuggestionsResult holds the list of autocomplete suggestion strings.
type SuggestionsResult struct {
	Suggestion []string `json:"suggestion"`
}

type foodsAutocompleteEnvelope struct {
	Suggestions foodsAutocompleteSuggestionsRaw `json:"suggestions"`
}

type foodsAutocompleteSuggestionsRaw struct {
	Suggestion json.RawMessage `json:"suggestion"`
}

// FoodsAutocomplete calls the food/autocomplete/v2 endpoint and returns search suggestions.
// Docs: https://platform.fatsecret.com/docs/v2/foods.autocomplete
func (c *FSNativeAPIClient) FoodsAutocomplete(params FoodsAutocompleteParams) (*FoodsAutocompleteResponse, error) {
	p := map[string]string{
		"format":     "json",
		"expression": params.Expression,
	}
	if params.MaxResults != nil {
		p["max_results"] = strconv.Itoa(*params.MaxResults)
	}
	if params.Region != nil {
		p["region"] = *params.Region
	}

	body, err := c.get(foodsAutocompletePath, p)
	if err != nil {
		return nil, err
	}

	var raw foodsAutocompleteEnvelope
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("fsnativeapi: unmarshal foods.autocomplete response: %w", err)
	}

	result := FoodsAutocompleteResponse{}

	// FatSecret returns `suggestion` as a string when there is exactly one result
	// and as an array otherwise.
	if len(raw.Suggestions.Suggestion) > 0 {
		switch raw.Suggestions.Suggestion[0] {
		case '[':
			if err := json.Unmarshal(raw.Suggestions.Suggestion, &result.Suggestions.Suggestion); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal suggestion array: %w", err)
			}
		case '"':
			var single string
			if err := json.Unmarshal(raw.Suggestions.Suggestion, &single); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal suggestion string: %w", err)
			}
			result.Suggestions.Suggestion = []string{single}
		}
	}

	return &result, nil
}