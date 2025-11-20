package main

import "net/http"
import "encoding/json"

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON( r *http.Request, data any) error {
	decoder:=json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	data:=map[string]string{
		"error":message,
	}	
	writeJSON(w,status,data)
}