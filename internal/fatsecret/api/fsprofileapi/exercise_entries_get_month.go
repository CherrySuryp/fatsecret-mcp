package fsprofileapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// https://platform.fatsecret.com/docs/v2/exercise_entries.get_month
const exerciseEntriesGetMonthPath = "exercise-entries/month/v2"

// ExerciseEntriesGetMonthResponse is the response from the exercise_entries.get_month endpoint.
type ExerciseEntriesGetMonthResponse struct {
	Month ExerciseEntriesMonth `json:"month"`
}

// ExerciseEntriesMonth holds the monthly summary of exercise entries.
type ExerciseEntriesMonth struct {
	FromDateInt string             `json:"from_date_int"`
	ToDateInt   string             `json:"to_date_int"`
	Days        []ExerciseEntryDay `json:"day"`
}

// ExerciseEntryDay is a daily calorie summary within a monthly exercise report.
type ExerciseEntryDay struct {
	DateInt  string `json:"date_int"`
	Calories string `json:"calories"`
}

// ExerciseEntriesGetMonth returns a monthly summary of exercise entries for the
// authenticated user. date is any day within the desired month (days since
// 1970-01-01); pass nil to use today.
//
// https://platform.fatsecret.com/docs/v2/exercise_entries.get_month
func (c *FSProfileClient) ExerciseEntriesGetMonth(date *int) (*ExerciseEntriesGetMonthResponse, error) {
	params := map[string]string{"format": "json"}
	if date != nil {
		params["date"] = strconv.Itoa(*date)
	}

	data, err := c.get(exerciseEntriesGetMonthPath, params)
	if err != nil {
		return nil, fmt.Errorf("ExerciseEntriesGetMonth: %w", err)
	}

	var resp ExerciseEntriesGetMonthResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("ExerciseEntriesGetMonth: decode response: %w", err)
	}

	return &resp, nil
}
