package fsprofileapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// https://platform.fatsecret.com/docs/v2/exercise_entries.get
const exerciseEntriesGetPath = "exercise-entries/v2"

// ExerciseEntriesGetResponse is the response from the exercise_entries.get endpoint.
type ExerciseEntriesGetResponse struct {
	ExerciseEntries ExerciseEntriesResult `json:"exercise_entries"`
}

// ExerciseEntriesResult holds the list of exercise entries.
type ExerciseEntriesResult struct {
	Entries []ExerciseEntry `json:"exercise_entry"`
}

// ExerciseEntry is a single exercise diary entry.
type ExerciseEntry struct {
	ExerciseID      string `json:"exercise_id"`
	ExerciseName    string `json:"exercise_name"`
	Minutes         string `json:"minutes"`
	Calories        string `json:"calories"`
	IsTemplateValue string `json:"is_template_value"`
}

// ExerciseEntriesGet returns exercise diary entries for the authenticated user.
// date is days since 1970-01-01; pass nil to use today.
//
// https://platform.fatsecret.com/docs/v2/exercise_entries.get
func (c *FSProfileClient) ExerciseEntriesGet(date *int) (*ExerciseEntriesGetResponse, error) {
	params := map[string]string{"format": "json"}
	if date != nil {
		params["date"] = strconv.Itoa(*date)
	}

	data, err := c.get(exerciseEntriesGetPath, params)
	if err != nil {
		return nil, fmt.Errorf("ExerciseEntriesGet: %w", err)
	}

	var resp ExerciseEntriesGetResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("ExerciseEntriesGet: decode response: %w", err)
	}

	return &resp, nil
}
