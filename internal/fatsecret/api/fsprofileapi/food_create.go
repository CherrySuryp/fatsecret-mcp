package fsprofileapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// https://platform.fatsecret.com/docs/v2/food.create
const foodCreatePath = "food/v2"

// FoodCreateParams holds all parameters for the food.create endpoint.
// Required fields: BrandType, BrandName, FoodName, ServingSize, Calories, Fat,
// Carbohydrate, Protein.
type FoodCreateParams struct {
	// Required
	BrandType    string  // "manufacturer", "restaurant", or "supermarket"
	BrandName    string  // brand identifier (e.g. "Quaker")
	FoodName     string  // food description, excluding the brand name
	ServingSize  string  // human-readable serving size description (e.g. "1 cup")
	Calories     float64 // energy in kcal
	Fat          float64 // total fat in grams
	Carbohydrate float64 // total carbohydrates in grams
	Protein      float64 // protein in grams

	// Optional — serving
	ServingAmount     string // quantity value for the serving
	ServingAmountUnit string // "g", "ml", or "oz" (default "g")

	// Optional — macros / fat breakdown
	CaloriesFromFat    *float64 // energy from fat in kcal
	SaturatedFat       *float64 // grams
	PolyunsaturatedFat *float64 // grams
	MonounsaturatedFat *float64 // grams
	TransFat           *float64 // grams

	// Optional — other nutrients
	Cholesterol *float64 // milligrams
	Sodium      *float64 // milligrams
	Potassium   *float64 // milligrams
	Fiber       *float64 // grams
	Sugar       *float64 // grams
	AddedSugars *float64 // grams

	// Optional — vitamins & minerals
	VitaminD *float64 // micrograms
	VitaminA *float64 // micrograms
	VitaminC *float64 // milligrams
	Calcium  *float64 // milligrams
	Iron     *float64 // milligrams

	// Optional — locale
	Region   string
	Language string
}

// FoodCreateResponse is the response from the food.create endpoint.
type FoodCreateResponse struct {
	FoodID struct {
		Value string `json:"value"`
	} `json:"food_id"`
}

// FoodCreate creates a new branded food entry for the authenticated user.
// Returns the newly created food_id on success.
//
// https://platform.fatsecret.com/docs/v2/food.create
func (c *FSProfileClient) FoodCreate(p FoodCreateParams) (*FoodCreateResponse, error) {
	params := map[string]string{
		"format":       "json",
		"brand_type":   p.BrandType,
		"brand_name":   p.BrandName,
		"food_name":    p.FoodName,
		"serving_size": p.ServingSize,
		"calories":     strconv.FormatFloat(p.Calories, 'f', -1, 64),
		"fat":          strconv.FormatFloat(p.Fat, 'f', -1, 64),
		"carbohydrate": strconv.FormatFloat(p.Carbohydrate, 'f', -1, 64),
		"protein":      strconv.FormatFloat(p.Protein, 'f', -1, 64),
	}

	if p.ServingAmount != "" {
		params["serving_amount"] = p.ServingAmount
	}
	if p.ServingAmountUnit != "" {
		params["serving_amount_unit"] = p.ServingAmountUnit
	}
	if p.CaloriesFromFat != nil {
		params["calories_from_fat"] = strconv.FormatFloat(*p.CaloriesFromFat, 'f', -1, 64)
	}
	if p.SaturatedFat != nil {
		params["saturated_fat"] = strconv.FormatFloat(*p.SaturatedFat, 'f', -1, 64)
	}
	if p.PolyunsaturatedFat != nil {
		params["polyunsaturated_fat"] = strconv.FormatFloat(*p.PolyunsaturatedFat, 'f', -1, 64)
	}
	if p.MonounsaturatedFat != nil {
		params["monounsaturated_fat"] = strconv.FormatFloat(*p.MonounsaturatedFat, 'f', -1, 64)
	}
	if p.TransFat != nil {
		params["trans_fat"] = strconv.FormatFloat(*p.TransFat, 'f', -1, 64)
	}
	if p.Cholesterol != nil {
		params["cholesterol"] = strconv.FormatFloat(*p.Cholesterol, 'f', -1, 64)
	}
	if p.Sodium != nil {
		params["sodium"] = strconv.FormatFloat(*p.Sodium, 'f', -1, 64)
	}
	if p.Potassium != nil {
		params["potassium"] = strconv.FormatFloat(*p.Potassium, 'f', -1, 64)
	}
	if p.Fiber != nil {
		params["fiber"] = strconv.FormatFloat(*p.Fiber, 'f', -1, 64)
	}
	if p.Sugar != nil {
		params["sugar"] = strconv.FormatFloat(*p.Sugar, 'f', -1, 64)
	}
	if p.AddedSugars != nil {
		params["added_sugars"] = strconv.FormatFloat(*p.AddedSugars, 'f', -1, 64)
	}
	if p.VitaminD != nil {
		params["vitamin_d"] = strconv.FormatFloat(*p.VitaminD, 'f', -1, 64)
	}
	if p.VitaminA != nil {
		params["vitamin_a"] = strconv.FormatFloat(*p.VitaminA, 'f', -1, 64)
	}
	if p.VitaminC != nil {
		params["vitamin_c"] = strconv.FormatFloat(*p.VitaminC, 'f', -1, 64)
	}
	if p.Calcium != nil {
		params["calcium"] = strconv.FormatFloat(*p.Calcium, 'f', -1, 64)
	}
	if p.Iron != nil {
		params["iron"] = strconv.FormatFloat(*p.Iron, 'f', -1, 64)
	}
	if p.Region != "" {
		params["region"] = p.Region
	}
	if p.Language != "" {
		params["language"] = p.Language
	}

	data, err := c.post(foodCreatePath, params)
	if err != nil {
		return nil, fmt.Errorf("FoodCreate: %w", err)
	}

	var resp FoodCreateResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("FoodCreate: decode response: %w", err)
	}

	return &resp, nil
}

