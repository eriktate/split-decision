package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func parseJSON[T any](r *http.Request) (T, error) {
	var parsed T
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		return parsed, fmt.Errorf("failed to read body: %w")
	}

	if err := json.Unmarshal(body, &parsed); err != nil {
		return parsed, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return parsed, nil
}
