package fsprofileapi

import (
	"encoding/json"
	"fmt"
)

// https://platform.fatsecret.com/docs/v2/exercises.get
const exercisesGetPath = "exercises/v2"

// ExercisesGetParams holds the optional parameters for the exercises.get endpoint.
type ExercisesGetParams struct {
	Region   string // optional
	Language string // optional; requires Region to be set
}

// ExercisesGetResponse is the response from the exercises.get endpoint.
type ExercisesGetResponse struct {
	ExerciseTypes ExerciseTypesResult `json:"exercise_types"`
}

// ExerciseTypesResult holds the list of exercise types.
type ExerciseTypesResult struct {
	Exercises []ExerciseType `json:"exercise"`
}

// ExerciseType is a single exercise type entry.
type ExerciseType struct {
	ExerciseID   string `json:"exercise_id"`
	ExerciseName string `json:"exercise_name"`
}

// ExercisesGet returns the list of exercise types available in FatSecret,
// optionally filtered by region and language.
//
// https://platform.fatsecret.com/docs/v2/exercises.get
func (c *FSProfileClient) ExercisesGet(p ExercisesGetParams) (*ExercisesGetResponse, error) {
	params := map[string]string{"format": "json"}
	if p.Region != "" {
		params["region"] = p.Region
	}
	if p.Language != "" {
		params["language"] = p.Language
	}

	data, err := c.get(exercisesGetPath, params)
	if err != nil {
		return nil, fmt.Errorf("ExercisesGet: %w", err)
	}

	var resp ExercisesGetResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("ExercisesGet: decode response: %w", err)
	}

	return &resp, nil
}
