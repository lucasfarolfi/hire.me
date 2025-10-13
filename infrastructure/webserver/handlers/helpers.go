package handlers

type HttpResponseErrorBody struct {
	ErrCode     string `json:"err_code"`
	Description string `json:"description"`
	Alias       string `json:"alias,omitempty"`
	URL         string `json:"url,omitempty"`
}
