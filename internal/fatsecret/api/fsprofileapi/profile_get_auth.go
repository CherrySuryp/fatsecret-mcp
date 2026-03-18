package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v1/profile.get_auth
const profileGetAuthPath = "profile/auth/v1"

// ProfileGetAuth retrieves the OAuth credentials (auth_token, auth_secret) for
// a profile. Pass userID to look up a profile by custom ID, or leave empty to
// use the profile associated with the current delegated auth.
//
// https://platform.fatsecret.com/docs/v1/profile.get_auth
func (c *FSProfileClient) ProfileGetAuth(userID string) (*ProfileAuthResponse, error) {
	params := map[string]string{"format": "json"}
	if userID != "" {
		params["user_id"] = userID
	}

	data, err := c.get(profileGetAuthPath, params)
	if err != nil {
		return nil, fmt.Errorf("ProfileGetAuth: %w", err)
	}

	var resp ProfileAuthResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("ProfileGetAuth: decode response: %w", err)
	}

	return &resp, nil
}
