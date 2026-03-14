package fsclient

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const weightUpdatePath = "weight/v1"

// WeightUpdateReq updates the user's weight for a given day.
// GoalWeightKg and CurrentHeightCm are required on the first weigh-in.
// Date must be within 2 days of today.
type WeightUpdateReq struct {
	CurrentWeightKg float64
	Date            *int
	WeightType      *WeightType
	HeightType      *HeightType
	GoalWeightKg    *float64
	CurrentHeightCm *float64
	Comment         *string
}

func (c *Client) WeightUpdate(payload WeightUpdateReq) (SuccessResp, error) {
	var result SuccessResp

	params := map[string]string{
		"current_weight_kg": fmt.Sprintf("%.3f", payload.CurrentWeightKg),
	}
	if payload.Date != nil {
		params["date"] = strconv.Itoa(*payload.Date)
	}
	if payload.WeightType != nil {
		params["weight_type"] = string(*payload.WeightType)
	}
	if payload.HeightType != nil {
		params["height_type"] = string(*payload.HeightType)
	}
	if payload.GoalWeightKg != nil {
		params["goal_weight_kg"] = fmt.Sprintf("%.3f", *payload.GoalWeightKg)
	}
	if payload.CurrentHeightCm != nil {
		params["current_height_cm"] = fmt.Sprintf("%.3f", *payload.CurrentHeightCm)
	}
	if payload.Comment != nil {
		params["comment"] = *payload.Comment
	}

	resp, err := c.post(weightUpdatePath, params)
	if err != nil {
		return result, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(data, &result)
	return result, err
}