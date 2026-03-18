package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v1/profile.get
const profileGetPath = "profile/v1"

// ProfileGetResponse is the response from the profile.get endpoint.
type ProfileGetResponse struct {
	Profile ProfileData `json:"profile"`
}

// ProfileData holds the user's profile information.
type ProfileData struct {
	GoalWeightKg      string `json:"goal_weight_kg"`
	HeightCm          string `json:"height_cm"`
	HeightMeasure     string `json:"height_measure"`
	LastWeightComment string `json:"last_weight_comment"`
	LastWeightDateInt string `json:"last_weight_date_int"`
	LastWeightKg      string `json:"last_weight_kg"`
	WeightMeasure     string `json:"weight_measure"`
}

// ProfileGet returns the profile information for the authenticated user.
//
// https://platform.fatsecret.com/docs/v1/profile.get
func (c *FSProfileClient) ProfileGet() (*ProfileGetResponse, error) {
	params := map[string]string{"format": "json"}

	data, err := c.get(profileGetPath, params)
	if err != nil {
		return nil, fmt.Errorf("ProfileGet: %w", err)
	}

	var resp ProfileGetResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("ProfileGet: decode response: %w", err)
	}

	return &resp, nil
}
