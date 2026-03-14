package fsclient

import (
	"encoding/json"
	"strconv"
)

const exerciseEntriesCommitDayPath = "exercise-entries/day/v1"

type ExerciseEntriesCommitDayReq struct {
	Date *int // optional; days since 1970-01-01
}

func (c *Client) ExerciseEntriesCommitDay(payload ExerciseEntriesCommitDayReq) (SuccessResp, error) {
	var result SuccessResp

	params := map[string]string{}
	if payload.Date != nil {
		params["date"] = strconv.Itoa(*payload.Date)
	}

	resp, err := c.post(exerciseEntriesCommitDayPath, params)
	if err != nil {
		return result, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(data, &result)
	return result, err
}