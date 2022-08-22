package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type ErrorResponse struct {
	Code   int    `json:"code"`
	Reason string `json:"reason,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	ContentTypeJSON = "application/json"
	ContentTypeText = "text/plain"
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// ServeJSON is a utility function to serve an arbitary object as JSON
func ServeJSON(w http.ResponseWriter, v interface{}, code, indent uint) error {
	w.Header().Add("Content-Type", ContentTypeJSON)
	w.WriteHeader(int(code))
	enc := json.NewEncoder(w)
	if indent > 0 {
		enc.SetIndent("", strings.Repeat(" ", int(indent)))
	}
	return enc.Encode(v)
}

// ServeText is a utility function to serve plaintext
func ServeText(w http.ResponseWriter, v string, code int) {
	w.Header().Add("Content-Type", ContentTypeText)
	w.WriteHeader(code)
	w.Write([]byte(v + "\n"))
}

// ServeError is a utility function to serve a JSON error notice
func ServeError(w http.ResponseWriter, code int, reason ...string) error {
	err := ErrorResponse{code, strings.Join(reason, " ")}
	if len(reason) == 0 {
		err.Reason = http.StatusText(code)
	}
	return ServeJSON(w, err, uint(code), 0)
}
