package models

import (
	"encoding/json"
	"testing"
)

type TestStruct struct {
	Date DateOnly `json:"date"`
}

func TestDateOnly_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name:    "DateOnly format",
			json:    `{"date": "2023-01-01"}`,
			wantErr: false,
		},
		{
			name:    "ISO8601 format",
			json:    `{"date": "2023-01-01T10:00:00Z"}`,
			wantErr: false,
		},
		{
			name:    "Invalid format",
			json:    `{"date": "invalid"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts TestStruct
			err := json.Unmarshal([]byte(tt.json), &ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				// We just want to ensure it parses without error for now.
				if ts.Date.Time.IsZero() {
					t.Errorf("Time should not be zero")
				}
			}
		})
	}
}
