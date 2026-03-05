package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/padiazg/pantry/internal/core/domain"
)

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data) //nolint:errcheck
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}

func respondHTTPError(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrNotFound) {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondError(w, http.StatusInternalServerError, err.Error())
}

// parseBoolQuery parses an optional boolean query parameter.
// Returns nil when the parameter is absent or unparseable.
func parseBoolQuery(r *http.Request, key string) *bool {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return nil
	}
	return &b
}
