package fsprofileapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// https://platform.fatsecret.com/docs/v1/exercise_entries.commit_day
const exerciseEntriesCommitDayPath = "exercise-entries/day/v1"

// ExerciseEntriesCommitDay commits exercise entries for a given day, making
// them permanent. date is days since 1970-01-01; pass nil to use today.
//
// https://platform.fatsecret.com/docs/v1/exercise_entries.commit_day
func (c *FSProfileClient) ExerciseEntriesCommitDay(date *int) (*SuccessResponse, error) {
	params := map[string]string{"format": "json"}
	if date != nil {
		params["date"] = strconv.Itoa(*date)
	}

	data, err := c.post(exerciseEntriesCommitDayPath, params)
	if err != nil {
		return nil, fmt.Errorf("ExerciseEntriesCommitDay: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("ExerciseEntriesCommitDay: decode response: %w", err)
	}

	return &resp, nil
}
