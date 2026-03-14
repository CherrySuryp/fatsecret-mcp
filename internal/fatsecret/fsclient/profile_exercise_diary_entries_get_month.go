package fsclient

import (
	"encoding/json"
	"strconv"
)

const exerciseEntriesGetMonthPath = "exercise-entries/month/v2"

type ExerciseEntriesGetMonthReq struct {
	Date *int // optional; any day within the target month (days since 1970-01-01)
}

type ExerciseEntriesDayResp struct {
	DateInt  int     `json:"date_int,string"`
	Calories float64 `json:"calories,string"`
}

type ExerciseEntriesGetMonthResp struct {
	FromDateInt int                      `json:"from_date_int,string"`
	ToDateInt   int                      `json:"to_date_int,string"`
	Day         []ExerciseEntriesDayResp `json:"day"`
}

type exerciseEntriesGetMonthWrapper struct {
	Month ExerciseEntriesGetMonthResp `json:"month"`
}

func (c *Client) ExerciseEntriesGetMonth(payload ExerciseEntriesGetMonthReq) (ExerciseEntriesGetMonthResp, error) {
	var result ExerciseEntriesGetMonthResp

	params := map[string]string{}
	if payload.Date != nil {
		params["date"] = strconv.Itoa(*payload.Date)
	}

	resp, err := c.get(exerciseEntriesGetMonthPath, params)
	if err != nil {
		return result, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return result, err
	}

	var wrapper exerciseEntriesGetMonthWrapper
	err = json.Unmarshal(data, &wrapper)
	return wrapper.Month, err
}