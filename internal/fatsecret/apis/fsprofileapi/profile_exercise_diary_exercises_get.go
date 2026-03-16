package fsprofileapi

import "encoding/json"

const exercisesGetPath = "exercises/v2"

type ExercisesGetReq struct {
	Region   *string // optional, Premier Exclusive
	Language *string // optional, Premier Exclusive
}

type ExercisesGetResp struct {
	ExerciseID   string `json:"exercise_id"`
	ExerciseName string `json:"exercise_name"`
}

type exercisesGetWrapper struct {
	ExerciseTypes struct {
		Exercise []ExercisesGetResp `json:"exercise"`
	} `json:"exercise_types"`
}

func (c *Client) ExercisesGet(payload ExercisesGetReq) ([]ExercisesGetResp, error) {
	params := map[string]string{}
	if payload.Region != nil {
		params["region"] = *payload.Region
	}
	if payload.Language != nil {
		params["language"] = *payload.Language
	}

	resp, err := c.get(exercisesGetPath, params)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var wrapper exercisesGetWrapper
	err = json.Unmarshal(data, &wrapper)
	return wrapper.ExerciseTypes.Exercise, err
}
