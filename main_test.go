package main

import (
	"testing"
)

func TestGenerateBuildNumber(t *testing.T) {
	tests := []struct {
		name string
		want uint
	}{
		// TODO: Add test cases.
		{"testing generation of build numbers",
			100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateBuildNumber(); got != tt.want {
				t.Errorf("GenerateBuildNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
