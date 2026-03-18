package fsnativeapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const foodGetPath = "food/v5"

// FoodGetParams holds the parameters for the food/v5 endpoint.
type FoodGetParams struct {
	// FoodID is the unique food identifier. Required.
	FoodID int64
	// IncludeSubCategories includes sub-category names in the response.
	IncludeSubCategories *bool
	// IncludeFoodImages includes food images. Requires premier.
	IncludeFoodImages *bool
	// IncludeFoodAttributes includes allergen and dietary preference data. Requires premier.
	IncludeFoodAttributes *bool
	// FlagDefaultServing marks the default serving in the response. Requires premier.
	FlagDefaultServing *bool
	// Region filters results by region (e.g. "US", "FR"). Defaults to "US".
	Region *string
	// Language specifies the response language. Requires Region to be set.
	Language *string
}

// FoodGetResponse is the top-level response envelope for food/v5.
type FoodGetResponse struct {
	Food FoodGetItem `json:"food"`
}

// FoodGetItem is a detailed food entry returned by food/v5.
type FoodGetItem struct {
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

// SubCategoryList wraps the sub-category string array.
type SubCategoryList struct {
	FoodSubCategory []string `json:"food_sub_category"`
}

// FoodImageList wraps the food image array.
type FoodImageList struct {
	FoodImage []FoodImage `json:"food_image"`
}

type foodGetEnvelope struct {
	Food foodGetItemRaw `json:"food"`
}

type foodGetItemRaw struct {
	FoodID            string          `json:"food_id"`
	FoodName          string          `json:"food_name"`
	FoodType          string          `json:"food_type"`
	FoodURL           string          `json:"food_url"`
	BrandName         string          `json:"brand_name,omitempty"`
	FoodSubCategories *SubCategoryList `json:"food_sub_categories,omitempty"`
	FoodImages        *FoodImageList   `json:"food_images,omitempty"`
	FoodAttributes    *FoodAttributes  `json:"food_attributes,omitempty"`
	Servings          struct {
		Serving json.RawMessage `json:"serving"`
	} `json:"servings"`
}

// FoodGet calls the food/v5 endpoint and returns detailed information for a single food item.
// Docs: https://platform.fatsecret.com/docs/v5/food.get
func (c *FSNativeAPIClient) FoodGet(params FoodGetParams) (*FoodGetResponse, error) {
	p := map[string]string{
		"format":  "json",
		"food_id": strconv.FormatInt(params.FoodID, 10),
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

	body, err := c.get(foodGetPath, p)
	if err != nil {
		return nil, err
	}

	var raw foodGetEnvelope
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("fsnativeapi: unmarshal food.get response: %w", err)
	}

	result := FoodGetResponse{
		Food: FoodGetItem{
			FoodID:            raw.Food.FoodID,
			FoodName:          raw.Food.FoodName,
			FoodType:          raw.Food.FoodType,
			FoodURL:           raw.Food.FoodURL,
			BrandName:         raw.Food.BrandName,
			FoodSubCategories: raw.Food.FoodSubCategories,
			FoodImages:        raw.Food.FoodImages,
			FoodAttributes:    raw.Food.FoodAttributes,
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