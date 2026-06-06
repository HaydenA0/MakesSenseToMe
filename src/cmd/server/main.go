package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"make/sense/to/me/src/converter"
)

type convertRequest struct {
	Values []float64 `json:"values"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("write json: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

func makeConverter(units []converter.Unit, mode string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req convertRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid json: %v", err))
			return
		}

		if len(req.Values) == 0 {
			writeError(w, http.StatusBadRequest, "values must not be empty")
			return
		}

		var results []converter.Result
		if mode == "absolute" {
			results = converter.ConvertAbsolute(req.Values, units)
		} else {
			results = converter.Convert(req.Values, units)
		}
		writeJSON(w, http.StatusOK, map[string]any{"results": results})
	}
}

func main() {
	addr := ":8080"

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/convert/time", makeConverter(converter.TimeUnits, "relative"))
	mux.HandleFunc("/convert/distance", makeConverter(converter.DistanceUnits, "relative"))
	mux.HandleFunc("/absolute/convert/time", makeConverter(converter.TimeUnits, "absolute"))
	mux.HandleFunc("/absolute/convert/distance", makeConverter(converter.DistanceUnits, "absolute"))

	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
