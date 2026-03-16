package fsprofileapi

import (
	"encoding/json"
	"strconv"
)

const exerciseEntriesTemplatePath = "exercise-entries/template/v1"

// ExerciseEntriesSaveTemplateReq saves exercise entries from the given date as
// template activities. Days is a bitmask where Sunday=bit1 … Saturday=bit7.
type ExerciseEntriesSaveTemplateReq struct {
	Days int  // required; bitmask (Sunday=1 … Saturday=64)
	Date *int // optional; days since 1970-01-01
}

func (c *Client) ExerciseEntriesSaveTemplate(payload ExerciseEntriesSaveTemplateReq) (SuccessResp, error) {
	var result SuccessResp

	params := map[string]string{
		"days": strconv.Itoa(payload.Days),
	}
	if payload.Date != nil {
		params["date"] = strconv.Itoa(*payload.Date)
	}

	resp, err := c.post(exerciseEntriesTemplatePath, params)
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
