package members_club

import (
	"encoding/json"
	"net/http"
)

type statusResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-type", "application/json")
	jsonError, _ := json.Marshal(errorResponse{message})
	http.Error(w, string(jsonError), statusCode)
}
