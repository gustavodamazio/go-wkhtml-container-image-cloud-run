package models

// RequestBody represents the expected structure of the incoming JSON request
type RequestBody struct {
	Data struct {
		HTML_STORAGE_PATH string `json:"html_storage_path"`
		CALLBACK_URL      string `json:"callback_url"`
		CALLBACK_METHOD   string `json:"callback_method"`
		CALLBACK_DATA     string `json:"callback_data"`
		CALLBACK_HEADERS  struct {
			Authorization string `json:"Authorization"`
		} `json:"callback_headers"`
	} `json:"data"`
}
