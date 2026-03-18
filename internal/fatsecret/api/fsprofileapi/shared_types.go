package fsprofileapi

// SuccessResponse is the common response shape returned by mutating endpoints
// that do not produce a new resource: {"success":{"value":"1"}}.
type SuccessResponse struct {
	Success struct {
		Value string `json:"value"`
	} `json:"success"`
}

