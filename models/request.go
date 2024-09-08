package models

// RequestBody represents the expected structure of the incoming JSON request
type RequestBody struct {
	Data struct {
		HTML string `json:"html"`
	} `json:"data"`
}
