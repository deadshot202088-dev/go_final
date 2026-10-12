package api

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := json.Marshal(data)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
