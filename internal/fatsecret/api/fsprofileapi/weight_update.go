package fsprofileapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// https://platform.fatsecret.com/docs/v1/weight.update
const weightUpdatePath = "weight/v1"

// WeightUpdateParams holds the parameters for the weight.update endpoint.
type WeightUpdateParams struct {
	CurrentWeightKg float64  // required; current weight in kg
	Date            *int     // optional; days since 1970-01-01, defaults to current day
	WeightType      string   // optional: "kg" or "lb" (default "kg")
	HeightType      string   // optional: "cm" or "inch" (default "cm")
	GoalWeightKg    *float64 // optional; goal weight in kg
	CurrentHeightCm *float64 // optional; required on first weigh-in
	Comment         string   // optional; comment for this weigh-in
}

// WeightUpdate records a weigh-in for the authenticated user.
//
// https://platform.fatsecret.com/docs/v1/weight.update
func (c *FSProfileClient) WeightUpdate(p WeightUpdateParams) (*SuccessResponse, error) {
	params := map[string]string{
		"format":            "json",
		"current_weight_kg": strconv.FormatFloat(p.CurrentWeightKg, 'f', -1, 64),
	}
	if p.Date != nil {
		params["date"] = strconv.Itoa(*p.Date)
	}
	if p.WeightType != "" {
		params["weight_type"] = p.WeightType
	}
	if p.HeightType != "" {
		params["height_type"] = p.HeightType
	}
	if p.GoalWeightKg != nil {
		params["goal_weight_kg"] = strconv.FormatFloat(*p.GoalWeightKg, 'f', -1, 64)
	}
	if p.CurrentHeightCm != nil {
		params["current_height_cm"] = strconv.FormatFloat(*p.CurrentHeightCm, 'f', -1, 64)
	}
	if p.Comment != "" {
		params["comment"] = p.Comment
	}

	data, err := c.post(weightUpdatePath, params)
	if err != nil {
		return nil, fmt.Errorf("WeightUpdate: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("WeightUpdate: decode response: %w", err)
	}

	return &resp, nil
}
