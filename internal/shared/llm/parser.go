package llm

import (
	"encoding/json"
	"fmt"
)

type TitlesResponse struct {
	Titles []string `json:"titles"`
}

func ParseTitlesJSON(raw string) (*TitlesResponse, error) {
	var parsed TitlesResponse
	if err := json.Unmarshal([]byte(raw), &parsed); err == nil {
		return &parsed, nil
	}

	trimmed := ExtractFirstJSONObject(raw)
	if trimmed == "" {
		return nil, fmt.Errorf("no json object found in response")
	}
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func ExtractFirstJSONObject(raw string) string {
	start := -1
	depth := 0
	for i, r := range raw {
		if r == '{' {
			if depth == 0 {
				start = i
			}
			depth++
			continue
		}
		if r == '}' {
			if depth > 0 {
				depth--
				if depth == 0 && start >= 0 {
					return raw[start : i+1]
				}
			}
		}
	}
	return ""
}
