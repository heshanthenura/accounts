package helpers

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	var body map[string]any

	if status >= 200 && status < 300 {
		body = map[string]any{
			"data": data,
		}
	} else {
		body = map[string]any{
			"error": map[string]any{
				"code":    status,
				"message": data,
			},
		}
	}

	json.NewEncoder(w).Encode(body)
}
