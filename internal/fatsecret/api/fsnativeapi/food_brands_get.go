package fsnativeapi

import (
	"encoding/json"
	"fmt"
)

const foodBrandsGetPath = "brands/v2"

// FoodBrandsGetParams holds the parameters for the brands/v2 endpoint.
// StartsWith is the primary filter; omit for popular brands, use "*" for numeric brands.
type FoodBrandsGetParams struct {
	// StartsWith filters brands by first letter. Use "*" for numeric brands.
	// Omit to retrieve popular brands.
	StartsWith string
	// BrandType filters by category: "manufacturer", "restaurant", or "supermarket".
	// Defaults to "manufacturer" when omitted.
	BrandType *string
	// Region filters results by region (e.g. "US", "FR"). Defaults to "US".
	Region *string
	// Language specifies the response language. Requires Region to be set.
	Language *string
}

// FoodBrandsGetResponse is the top-level response envelope for brands/v2.
type FoodBrandsGetResponse struct {
	FoodBrands FoodBrandsResult `json:"food_brands"`
}

// FoodBrandsResult holds the list of brand name strings.
type FoodBrandsResult struct {
	FoodBrand []string `json:"food_brand"`
}

type foodBrandsEnvelope struct {
	FoodBrands foodBrandsResultRaw `json:"food_brands"`
}

type foodBrandsResultRaw struct {
	FoodBrand json.RawMessage `json:"food_brand"`
}

// FoodBrandsGet calls the brands/v2 endpoint and returns matching food brand names.
// Docs: https://platform.fatsecret.com/docs/v2/food_brands.get
func (c *FSNativeAPIClient) FoodBrandsGet(params FoodBrandsGetParams) (*FoodBrandsGetResponse, error) {
	p := map[string]string{
		"format": "json",
	}
	if params.StartsWith != "" {
		p["starts_with"] = params.StartsWith
	}
	if params.BrandType != nil {
		p["brand_type"] = *params.BrandType
	}
	if params.Region != nil {
		p["region"] = *params.Region
	}
	if params.Language != nil {
		p["language"] = *params.Language
	}

	body, err := c.get(foodBrandsGetPath, p)
	if err != nil {
		return nil, err
	}

	var raw foodBrandsEnvelope
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("fsnativeapi: unmarshal food_brands.get response: %w", err)
	}

	result := FoodBrandsGetResponse{}

	// FatSecret returns `food_brand` as a string when there is exactly one result
	// and as an array otherwise.
	if len(raw.FoodBrands.FoodBrand) > 0 {
		switch raw.FoodBrands.FoodBrand[0] {
		case '[':
			if err := json.Unmarshal(raw.FoodBrands.FoodBrand, &result.FoodBrands.FoodBrand); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal food_brand array: %w", err)
			}
		case '"':
			var single string
			if err := json.Unmarshal(raw.FoodBrands.FoodBrand, &single); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal food_brand string: %w", err)
			}
			result.FoodBrands.FoodBrand = []string{single}
		}
	}

	return &result, nil
}