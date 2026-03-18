package fsnativeapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const foodFindIDForBarcodePath = "food/barcode/find-by-id/v2"

// FoodFindIDForBarcodeParams holds the parameters for the food/barcode/find-by-id/v2 endpoint.
type FoodFindIDForBarcodeParams struct {
	// Barcode is the 13-digit GTIN-13 barcode. Required.
	// UPC-A and EAN-8 are also accepted; UPC-E must be converted to UPC-A first.
	Barcode string
	// IncludeSubCategories includes sub-category names in the response.
	IncludeSubCategories *bool
	// IncludeFoodImages includes food images. Requires premier.
	IncludeFoodImages *bool
	// IncludeFoodAttributes includes allergen and dietary preference data. Requires premier.
	IncludeFoodAttributes *bool
	// FlagDefaultServing marks the default serving in the response.
	FlagDefaultServing *bool
	// Region filters results by region (e.g. "US", "FR"). Defaults to "US".
	Region *string
	// Language specifies the response language. Requires Region to be set.
	Language *string
}

// FoodFindIDForBarcodeResponse is the top-level response envelope for food/barcode/find-by-id/v2.
type FoodFindIDForBarcodeResponse struct {
	Food BarcodeFood `json:"food"`
}

// BarcodeFood is the food item matched to the scanned barcode.
type BarcodeFood struct {
	FoodID         string          `json:"food_id"`
	FoodName       string          `json:"food_name"`
	BrandName      string          `json:"brand_name,omitempty"`
	FoodType       string          `json:"food_type"`
	FoodURL        string          `json:"food_url"`
	FoodAttributes *FoodAttributes `json:"food_attributes,omitempty"`
	Servings       struct {
		Serving []Serving `json:"serving"`
	} `json:"servings"`
}

type foodFindIDForBarcodeEnvelope struct {
	Food barcodeFoodRaw `json:"food"`
}

type barcodeFoodRaw struct {
	FoodID         string          `json:"food_id"`
	FoodName       string          `json:"food_name"`
	BrandName      string          `json:"brand_name,omitempty"`
	FoodType       string          `json:"food_type"`
	FoodURL        string          `json:"food_url"`
	FoodAttributes *FoodAttributes `json:"food_attributes,omitempty"`
	Servings       struct {
		Serving json.RawMessage `json:"serving"`
	} `json:"servings"`
}

// FoodFindIDForBarcode calls the food/barcode/find-by-id/v2 endpoint and returns
// the food item matching the given barcode.
// Docs: https://platform.fatsecret.com/docs/v2/food.find_id_for_barcode
func (c *FSNativeAPIClient) FoodFindIDForBarcode(params FoodFindIDForBarcodeParams) (*FoodFindIDForBarcodeResponse, error) {
	p := map[string]string{
		"format":  "json",
		"barcode": params.Barcode,
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
	if params.Region != nil {
		p["region"] = *params.Region
	}
	if params.Language != nil {
		p["language"] = *params.Language
	}

	body, err := c.get(foodFindIDForBarcodePath, p)
	if err != nil {
		return nil, err
	}

	var raw foodFindIDForBarcodeEnvelope
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("fsnativeapi: unmarshal food.find_id_for_barcode response: %w", err)
	}

	result := FoodFindIDForBarcodeResponse{
		Food: BarcodeFood{
			FoodID:         raw.Food.FoodID,
			FoodName:       raw.Food.FoodName,
			BrandName:      raw.Food.BrandName,
			FoodType:       raw.Food.FoodType,
			FoodURL:        raw.Food.FoodURL,
			FoodAttributes: raw.Food.FoodAttributes,
		},
	}

	// FatSecret returns `serving` as an object when there is exactly one result
	// and as an array otherwise.
	if len(raw.Food.Servings.Serving) > 0 {
		switch raw.Food.Servings.Serving[0] {
		case '[':
			if err := json.Unmarshal(raw.Food.Servings.Serving, &result.Food.Servings.Serving); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal serving array: %w", err)
			}
		case '{':
			var single Serving
			if err := json.Unmarshal(raw.Food.Servings.Serving, &single); err != nil {
				return nil, fmt.Errorf("fsnativeapi: unmarshal serving object: %w", err)
			}
			result.Food.Servings.Serving = []Serving{single}
		}
	}

	return &result, nil
}