package fsprofileapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// https://platform.fatsecret.com/docs/v1/exercise_entry.edit
const exerciseEntryEditPath = "exercise-entries/template/v1"

// ExerciseEntryEditParams holds the parameters for the exercise_entry.edit endpoint.
// This is a shift operation: it replaces ShiftFromID with ShiftToID in the
// exercise diary for the given date.
type ExerciseEntryEditParams struct {
	ShiftToID   string // required; exercise type ID to shift to (use "0" for custom)
	ShiftFromID string // required; exercise type ID to shift from (use "0" for custom)
	Minutes     int    // required; duration in minutes

	Date          *int   // optional; days since 1970-01-01, defaults to current day
	ShiftToName   string // required when ShiftToID is "0"
	ShiftFromName string // required when ShiftFromID is "0"
	Kcal          *int   // required when ShiftToID is "0"; calories burned
}

// ExerciseEntryEdit replaces one exercise entry with another in the authenticated
// user's diary. Use ShiftToID/ShiftFromID of "0" with corresponding name and
// kcal fields for custom (unlisted) exercises.
//
// https://platform.fatsecret.com/docs/v1/exercise_entry.edit
func (c *FSProfileClient) ExerciseEntryEdit(p ExerciseEntryEditParams) (*SuccessResponse, error) {
	params := map[string]string{
		"format":        "json",
		"shift_to_id":   p.ShiftToID,
		"shift_from_id": p.ShiftFromID,
		"minutes":       strconv.Itoa(p.Minutes),
	}
	if p.Date != nil {
		params["date"] = strconv.Itoa(*p.Date)
	}
	if p.ShiftToName != "" {
		params["shift_to_name"] = p.ShiftToName
	}
	if p.ShiftFromName != "" {
		params["shift_from_name"] = p.ShiftFromName
	}
	if p.Kcal != nil {
		params["kcal"] = strconv.Itoa(*p.Kcal)
	}

	data, err := c.put(exerciseEntryEditPath, params)
	if err != nil {
		return nil, fmt.Errorf("ExerciseEntryEdit: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("ExerciseEntryEdit: decode response: %w", err)
	}

	return &resp, nil
}
