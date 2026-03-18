package fsnativeapi

// Serving represents a single serving size with full nutritional data.
// All numeric values are returned as strings by the FatSecret API.
type Serving struct {
	ServingID              string `json:"serving_id"`
	ServingDescription     string `json:"serving_description"`
	ServingURL             string `json:"serving_url,omitempty"`
	MetricServingAmount    string `json:"metric_serving_amount,omitempty"`
	MetricServingUnit      string `json:"metric_serving_unit,omitempty"`
	NumberOfUnits          string `json:"number_of_units,omitempty"`
	MeasurementDescription string `json:"measurement_description,omitempty"`
	IsDefault              string `json:"is_default,omitempty"`
	Calories               string `json:"calories,omitempty"`
	Carbohydrate           string `json:"carbohydrate,omitempty"`
	Protein                string `json:"protein,omitempty"`
	Fat                    string `json:"fat,omitempty"`
	SaturatedFat           string `json:"saturated_fat,omitempty"`
	PolyunsaturatedFat     string `json:"polyunsaturated_fat,omitempty"`
	MonounsaturatedFat     string `json:"monounsaturated_fat,omitempty"`
	TransFat               string `json:"trans_fat,omitempty"`
	Cholesterol            string `json:"cholesterol,omitempty"`
	Sodium                 string `json:"sodium,omitempty"`
	Potassium              string `json:"potassium,omitempty"`
	Fiber                  string `json:"fiber,omitempty"`
	Sugar                  string `json:"sugar,omitempty"`
	AddedSugars            string `json:"added_sugars,omitempty"`
	VitaminA               string `json:"vitamin_a,omitempty"`
	VitaminC               string `json:"vitamin_c,omitempty"`
	VitaminD               string `json:"vitamin_d,omitempty"`
	Calcium                string `json:"calcium,omitempty"`
	Iron                   string `json:"iron,omitempty"`
}

// FoodAttributeEntry is a single allergen or dietary preference flag.
// Value: "1" = yes/true, "0" = no/false, "-1" = unknown.
type FoodAttributeEntry struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

// FoodAttributes holds allergen and dietary preference data for a food item.
type FoodAttributes struct {
	Allergens   []FoodAttributeEntry `json:"allergens,omitempty"`
	Preferences []FoodAttributeEntry `json:"preferences,omitempty"`
}

// FoodImage is a single image associated with a food item.
type FoodImage struct {
	ImageURL  string `json:"image_url"`
	ImageType string `json:"image_type"` // "Standard" or "Isolated"
}

// EatenFood is a previously consumed food item provided for improved matching
// accuracy in NLP and image recognition requests.
type EatenFood struct {
	FoodID             int64  `json:"food_id"`
	FoodName           string `json:"food_name"`
	FoodBrand          string `json:"food_brand,omitempty"`
	ServingDescription string `json:"serving_description,omitempty"`
	ServingSize        string `json:"serving_size,omitempty"`
}

// NutritionalContent holds macronutrient and micronutrient totals returned by
// the NLP and image recognition APIs. All values are decimal strings.
type NutritionalContent struct {
	Calories           string `json:"calories,omitempty"`
	Carbohydrate       string `json:"carbohydrate,omitempty"`
	Protein            string `json:"protein,omitempty"`
	Fat                string `json:"fat,omitempty"`
	SaturatedFat       string `json:"saturated_fat,omitempty"`
	PolyunsaturatedFat string `json:"polyunsaturated_fat,omitempty"`
	MonounsaturatedFat string `json:"monounsaturated_fat,omitempty"`
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

// EatenDetails describes the quantity and nutritional summary of a detected food.
type EatenDetails struct {
	FoodNameSingular        string             `json:"food_name_singular,omitempty"`
	FoodNamePlural          string             `json:"food_name_plural,omitempty"`
	SingularDescription     string             `json:"singular_description,omitempty"`
	PluralDescription       string             `json:"plural_description,omitempty"`
	Units                   string             `json:"units,omitempty"`
	MetricDescription       string             `json:"metric_description,omitempty"`
	TotalMetricAmount       string             `json:"total_metric_amount,omitempty"`
	PerUnitMetricAmount     string             `json:"per_unit_metric_amount,omitempty"`
	TotalNutritionalContent NutritionalContent `json:"total_nutritional_content,omitempty"`
}

// SuggestedServing is the API's recommended serving for a detected food item.
type SuggestedServing struct {
	ServingID                string `json:"serving_id"`
	ServingDescription       string `json:"serving_description,omitempty"`
	CustomServingDescription string `json:"custom_serving_description,omitempty"`
	MetricServingDescription string `json:"metric_serving_description,omitempty"`
	MetricMeasureAmount      string `json:"metric_measure_amount,omitempty"`
	NumberOfUnits            string `json:"number_of_units,omitempty"`
}

// FoodResponseItem is a single detected food entry in NLP and image recognition responses.
type FoodResponseItem struct {
	FoodID          string           `json:"food_id"`
	FoodEntryName   string           `json:"food_entry_name"`
	Eaten           EatenDetails     `json:"eaten"`
	SuggestedServing SuggestedServing `json:"suggested_serving"`
	// Food is only populated when include_food_data is true.
	Food *AIFoodDetail `json:"food,omitempty"`
}

// AIFoodDetail is the optional full food data included in NLP and image recognition
// responses when include_food_data is true.
type AIFoodDetail struct {
	FoodID   string `json:"food_id"`
	FoodName string `json:"food_name"`
	FoodType string `json:"food_type"`
	FoodURL  string `json:"food_url"`
	Servings struct {
		Serving []Serving `json:"serving"`
	} `json:"servings"`
}