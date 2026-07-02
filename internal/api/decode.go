package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// decodeJSON reads and decodes a JSON request body into dst.
// Returns an error if the body is missing or malformed.
func decodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is required")
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}
