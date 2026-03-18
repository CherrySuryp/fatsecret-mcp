package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v1/profile.create
const profileCreatePath = "profile/v1"

// ProfileAuthResponse holds the OAuth credentials returned when creating a
// profile or retrieving profile auth. These tokens are used to construct a
// DelegatedAuth for subsequent user-scoped requests.
type ProfileAuthResponse struct {
	Profile struct {
		AuthToken  string `json:"auth_token"`
		AuthSecret string `json:"auth_secret"`
	} `json:"profile"`
}

// ProfileCreate creates a new FatSecret user profile with the given user_id.
// Returns OAuth credentials (auth_token, auth_secret) for the new profile.
//
// https://platform.fatsecret.com/docs/v1/profile.create
func (c *FSProfileClient) ProfileCreate(userID string) (*ProfileAuthResponse, error) {
	params := map[string]string{
		"format":  "json",
		"user_id": userID,
	}

	data, err := c.post(profileCreatePath, params)
	if err != nil {
		return nil, fmt.Errorf("ProfileCreate: %w", err)
	}

	var resp ProfileAuthResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("ProfileCreate: decode response: %w", err)
	}

	return &resp, nil
}
