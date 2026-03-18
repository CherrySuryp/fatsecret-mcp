package fsprofileapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// https://platform.fatsecret.com/docs/v1/exercise_entries.save_template
const exerciseEntriesSaveTemplatePath = "exercise-entries/template/v1"

// ExerciseEntriesSaveTemplateParams holds the parameters for the
// exercise_entries.save_template endpoint.
type ExerciseEntriesSaveTemplateParams struct {
	// Days is a bitmask of weekdays on which to apply the template.
	// Sunday = bit 0 (1), Monday = bit 1 (2), ..., Saturday = bit 6 (64).
	// Use 0–127 to select any combination of days.
	Days int  // required
	Date *int // optional; days since 1970-01-01, defaults to current day
}

// ExerciseEntriesSaveTemplate saves the current day's exercise entries as a
// weekly template for the specified days of the week.
//
// https://platform.fatsecret.com/docs/v1/exercise_entries.save_template
func (c *FSProfileClient) ExerciseEntriesSaveTemplate(p ExerciseEntriesSaveTemplateParams) (*SuccessResponse, error) {
	params := map[string]string{
		"format": "json",
		"days":   strconv.Itoa(p.Days),
	}
	if p.Date != nil {
		params["date"] = strconv.Itoa(*p.Date)
	}

	data, err := c.post(exerciseEntriesSaveTemplatePath, params)
	if err != nil {
		return nil, fmt.Errorf("ExerciseEntriesSaveTemplate: %w", err)
	}

	var resp SuccessResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("ExerciseEntriesSaveTemplate: decode response: %w", err)
	}

	return &resp, nil
}
