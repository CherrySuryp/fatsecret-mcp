package fsclient

type MealType string

const (
	MealTypeBreakfast MealType = "breakfast"
	MealTypeLunch     MealType = "lunch"
	MealTypeDinner    MealType = "dinner"
	MealTypeOther     MealType = "other"
)

type WeightType string

const (
	WeightTypeKg WeightType = "kg"
	WeightTypeLb WeightType = "lb"
)

type HeightType string

const (
	HeightTypeCm   HeightType = "cm"
	HeightTypeInch HeightType = "inch"
)

type SuccessResp struct {
	Success struct {
		Value int `json:"value,string"`
	} `json:"success"`
}

func (r SuccessResp) IsSuccess() bool {
	return r.Success.Value == 1
}
