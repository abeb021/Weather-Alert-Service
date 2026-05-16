package handlers

import (
	"encoding/json"
	"net/http"
)

/*
{
  "type": "https://your-api.com/errors/city-not-found",
  "title": "City Not Found",
  "status": 404,
  "detail": "The city 'Londoon' was not found in the weather database.",
  "instance": "/api/weather/current?city=Londoon"
}
*/
type ProblemDetail struct {
    Type     string `json:"type"`
    Title    string `json:"title"`
    Status   int    `json:"status"`
    Detail   string `json:"detail"`
    Instance string `json:"instance,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeProblem(w http.ResponseWriter, r *http.Request, status int, title, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ProblemDetail{
		Type: "about:blank",
		Title: title,
		Detail: detail,
		Instance: r.URL.String(),
	})
}