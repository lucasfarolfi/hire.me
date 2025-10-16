package handlers

import (
	"encoding/json"
	"net/http"
)

type HttpResponseErrorBody struct {
	ErrCode     string `json:"err_code"`
	Description string `json:"description"`
	Alias       string `json:"alias,omitempty"`
}

func retrieveErrorResponseBody(w http.ResponseWriter, statusCode int, errCode, description, alias string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(&HttpResponseErrorBody{errCode, description, alias}); err != nil {
		http.Error(w, "failed to encode error body", http.StatusInternalServerError)
	}
}
