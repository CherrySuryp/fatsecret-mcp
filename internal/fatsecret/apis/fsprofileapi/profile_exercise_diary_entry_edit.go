package fsprofileapi

import (
	"encoding/json"
	"strconv"
)

// ExerciseEntryEditReq shifts minutes from one exercise type to another.
// ShiftToName is required when ShiftToID is 0 (custom exercise).
// ShiftFromName is required when ShiftFromID is 0 (custom exercise).
// Kcal is required when ShiftToID is 0.
type ExerciseEntryEditReq struct {
	ShiftToID     int
	ShiftFromID   int
	Minutes       int
	Date          *int
	ShiftToName   *string
	ShiftFromName *string
	Kcal          *int
}

func (c *Client) ExerciseEntryEdit(payload ExerciseEntryEditReq) (SuccessResp, error) {
	var result SuccessResp

	params := map[string]string{
		"shift_to_id":   strconv.Itoa(payload.ShiftToID),
		"shift_from_id": strconv.Itoa(payload.ShiftFromID),
		"minutes":       strconv.Itoa(payload.Minutes),
	}
	if payload.Date != nil {
		params["date"] = strconv.Itoa(*payload.Date)
	}
	if payload.ShiftToName != nil {
		params["shift_to_name"] = *payload.ShiftToName
	}
	if payload.ShiftFromName != nil {
		params["shift_from_name"] = *payload.ShiftFromName
	}
	if payload.Kcal != nil {
		params["kcal"] = strconv.Itoa(*payload.Kcal)
	}

	resp, err := c.put(exerciseEntriesTemplatePath, params)
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
