package fsprofileapi

import (
	"encoding/json"
	"strconv"
)

const exerciseEntriesGetPath = "exercise-entries/v2"

type ExerciseEntriesGetReq struct {
	Date *int // optional; days since 1970-01-01
}

type ExerciseEntriesGetResp struct {
	ExerciseID      string  `json:"exercise_id"`
	ExerciseName    string  `json:"exercise_name"`
	IsTemplateValue int     `json:"is_template_value,string"`
	Minutes         int     `json:"minutes,string"`
	Calories        float64 `json:"calories,string"`
}

type exerciseEntriesGetWrapper struct {
	ExerciseEntries struct {
		ExerciseEntry []ExerciseEntriesGetResp `json:"exercise_entry"`
	} `json:"exercise_entries"`
}

func (c *Client) ExerciseEntriesGet(payload ExerciseEntriesGetReq) ([]ExerciseEntriesGetResp, error) {
	params := map[string]string{}
	if payload.Date != nil {
		params["date"] = strconv.Itoa(*payload.Date)
	}

	resp, err := c.get(exerciseEntriesGetPath, params)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var wrapper exerciseEntriesGetWrapper
	err = json.Unmarshal(data, &wrapper)
	return wrapper.ExerciseEntries.ExerciseEntry, err
}
