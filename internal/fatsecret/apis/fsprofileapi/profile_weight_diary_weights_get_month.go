package fsprofileapi

import (
	"encoding/json"
	"strconv"
)

const weightsGetMonthPath = "weight/month/v2"

type WeightsGetMonthReq struct {
	Date *int // optional; any day within the target month (days since 1970-01-01)
}

type WeightDayResp struct {
	DateInt       int     `json:"date_int,string"`
	WeightKg      float64 `json:"weight_kg,string"`
	WeightComment string  `json:"weight_comment"`
}

type WeightsGetMonthResp struct {
	FromDateInt int             `json:"from_date_int,string"`
	ToDateInt   int             `json:"to_date_int,string"`
	Day         []WeightDayResp `json:"day"`
}

type weightsGetMonthWrapper struct {
	Month WeightsGetMonthResp `json:"month"`
}

func (c *Client) WeightsGetMonth(payload WeightsGetMonthReq) (WeightsGetMonthResp, error) {
	var result WeightsGetMonthResp

	params := map[string]string{}
	if payload.Date != nil {
		params["date"] = strconv.Itoa(*payload.Date)
	}

	resp, err := c.get(weightsGetMonthPath, params)
	if err != nil {
		return result, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return result, err
	}

	var wrapper weightsGetMonthWrapper
	err = json.Unmarshal(data, &wrapper)
	return wrapper.Month, err
}
