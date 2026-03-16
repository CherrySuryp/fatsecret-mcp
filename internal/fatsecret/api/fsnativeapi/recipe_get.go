package fsnativeapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const recipeGetPath = "recipe/v2"

// RecipeGetParams holds the parameters for the recipe/v2 endpoint.
type RecipeGetParams struct {
	// RecipeID is the unique recipe identifier. Required.
	RecipeID int64
	// Region filters results by region (e.g. "US", "FR"). Defaults to "US". Premier only.
	Region *string
}

// RecipeGetResponse is the top-level response envelope for recipe/v2.
type RecipeGetResponse struct {
	Recipe RecipeItem `json:"recipe"`
}

// RecipeItem is a full recipe entry returned by recipe/v2.
type RecipeItem struct {
	RecipeID           string                `json:"recipe_id"`
	RecipeName         string                `json:"recipe_name"`
	RecipeURL          string                `json:"recipe_url"`
	RecipeDescription  string                `json:"recipe_description,omitempty"`
	NumberOfServings   string                `json:"number_of_servings"`
	GramsPerPortion    string                `json:"grams_per_portion,omitempty"`
	PreparationTimeMin string                `json:"preparation_time_min,omitempty"`
	CookingTimeMin     string                `json:"cooking_time_min,omitempty"`
	Rating             string                `json:"rating,omitempty"`
	RecipeTypes        *RecipeTypeList       `json:"recipe_types,omitempty"`
	RecipeCategories   *RecipeCategoryList   `json:"recipe_categories,omitempty"`
	RecipeImages       *RecipeImageList      `json:"recipe_images,omitempty"`
	ServingSizes       *RecipeServingSizes   `json:"serving_sizes,omitempty"`
	Ingredients        *RecipeIngredientList `json:"ingredients,omitempty"`
	Directions         *RecipeDirectionList  `json:"directions,omitempty"`
}

// RecipeTypeList wraps the recipe type string array.
type RecipeTypeList struct {
	RecipeType []string `json:"recipe_type"`
}

// RecipeCategoryList wraps the recipe category array.
type RecipeCategoryList struct {
	RecipeCategory []RecipeCategory `json:"recipe_category"`
}

// RecipeCategory is a single recipe category with a name and URL.
type RecipeCategory struct {
	RecipeCategoryName string `json:"recipe_category_name"`
	RecipeCategoryURL  string `json:"recipe_category_url"`
}

// RecipeImageList wraps the recipe image URL array.
type RecipeImageList struct {
	RecipeImage []string `json:"recipe_image"`
}

// RecipeServingSizes wraps the serving size array.
type RecipeServingSizes struct {
	Serving []RecipeServing `json:"serving"`
}

// RecipeServing holds nutritional data for one serving of a recipe.
type RecipeServing struct {
	ServingSize        string `json:"serving_size"`
	Calories           string `json:"calories,omitempty"`
	Carbohydrate       string `json:"carbohydrate,omitempty"`
	Protein            string `json:"protein,omitempty"`
	Fat                string `json:"fat,omitempty"`
	SaturatedFat       string `json:"saturated_fat,omitempty"`
	PolyunsaturatedFat string `json:"polyunsaturated_fat,omitempty"`
	MonounsaturatedFat string `json:"monounsaturated_fat,omitempty"`
	TransFat           string `json:"trans_fat,omitempty"`
	Cholesterol        string `json:"cholesterol,omitempty"`
	Sodium             string `json:"sodium,omitempty"`
	Potassium          string `json:"potassium,omitempty"`
	Fiber              string `json:"fiber,omitempty"`
	Sugar              string `json:"sugar,omitempty"`
	VitaminA           string `json:"vitamin_a,omitempty"`
	VitaminC           string `json:"vitamin_c,omitempty"`
	Calcium            string `json:"calcium,omitempty"`
	Iron               string `json:"iron,omitempty"`
}

// RecipeIngredientList wraps the ingredient array.
type RecipeIngredientList struct {
	Ingredient []RecipeIngredient `json:"ingredient"`
}

// RecipeIngredient is a single ingredient in a recipe.
type RecipeIngredient struct {
	FoodID                 string `json:"food_id"`
	ServingID              string `json:"serving_id"`
	FoodName               string `json:"food_name"`
	MeasurementDescription string `json:"measurement_description"`
	NumberOfUnits          string `json:"number_of_units"`
	IngredientURL          string `json:"ingredient_url"`
	IngredientDescription  string `json:"ingredient_description"`
}

// RecipeDirectionList wraps the direction step array.
type RecipeDirectionList struct {
	Direction []RecipeDirection `json:"direction"`
}

// RecipeDirection is a single preparation step in a recipe.
type RecipeDirection struct {
	DirectionNumber      string `json:"direction_number"`
	DirectionDescription string `json:"direction_description"`
}

// RecipeGet calls the recipe/v2 endpoint and returns the full recipe for the given ID.
// Docs: https://platform.fatsecret.com/docs/v2/recipe.get
func (c *FSNativeAPIClient) RecipeGet(params RecipeGetParams) (*RecipeGetResponse, error) {
	p := map[string]string{
		"format":    "json",
		"recipe_id": strconv.FormatInt(params.RecipeID, 10),
	}
	if params.Region != nil {
		p["region"] = *params.Region
	}

	body, err := c.get(recipeGetPath, p)
	if err != nil {
		return nil, err
	}

	var result RecipeGetResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("fsnativeapi: unmarshal recipe.get response: %w", err)
	}

	return &result, nil
}
