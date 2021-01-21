package router

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, data interface{}, code int) {
	buf, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	status := http.StatusOK
	if code != 0 {
		status = code
	}
	w.WriteHeader(status)
	w.Write(buf)
}
