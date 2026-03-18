package fsprofileapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// https://platform.fatsecret.com/docs/v2/weights.get_month
const weightsGetMonthPath = "weight/month/v2"

// WeightsGetMonthResponse is the response from the weights.get_month endpoint.
type WeightsGetMonthResponse struct {
	Month WeightsMonth `json:"month"`
}

// WeightsMonth holds the monthly summary of weigh-ins.
type WeightsMonth struct {
	FromDateInt string      `json:"from_date_int"`
	ToDateInt   string      `json:"to_date_int"`
	Days        []WeightDay `json:"day"`
}

// WeightDay is a single weigh-in entry within a monthly report.
type WeightDay struct {
	DateInt       string `json:"date_int"`
	WeightKg      string `json:"weight_kg"`
	WeightComment string `json:"weight_comment"`
}

// WeightsGetMonth returns a monthly summary of weigh-ins for the authenticated
// user. date is any day within the desired month (days since 1970-01-01);
// pass nil to use today.
//
// https://platform.fatsecret.com/docs/v2/weights.get_month
func (c *FSProfileClient) WeightsGetMonth(date *int) (*WeightsGetMonthResponse, error) {
	params := map[string]string{"format": "json"}
	if date != nil {
		params["date"] = strconv.Itoa(*date)
	}

	data, err := c.get(weightsGetMonthPath, params)
	if err != nil {
		return nil, fmt.Errorf("WeightsGetMonth: %w", err)
	}

	var resp WeightsGetMonthResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("WeightsGetMonth: decode response: %w", err)
	}

	return &resp, nil
}
