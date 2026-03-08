package main

import "testing"

func TestValidateFlags(t *testing.T) {
	tests := []struct {
		name        string
		target      string
		startPort   int
		endPort     int
		top         int
		timeout     int
		concurrency int
		wantErr     bool
	}{
		{"valid range", "localhost", 1, 1024, 0, 300, 1000, false},
		{"valid top", "localhost", 0, 0, 100, 300, 1000, false},
		{"missing target", "", 1, 1024, 0, 300, 1000, true},
		{"start port too low", "localhost", 0, 1024, 0, 300, 1000, true},
		{"start port too high", "localhost", 65536, 1024, 0, 300, 1000, true},
		{"end port too low", "localhost", 1, 0, 0, 300, 1000, true},
		{"end port too high", "localhost", 1, 65536, 0, 300, 1000, true},
		{"start greater than end", "localhost", 100, 50, 0, 300, 1000, true},
		{"top too low", "localhost", 0, 0, 0, 300, 1000, true},
		{"top too high", "localhost", 0, 0, 99999, 300, 1000, true},
		{"timeout too low", "localhost", 1, 1024, 0, 0, 1000, true},
		{"concurrency too low", "localhost", 1, 1024, 0, 300, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFlags(tt.target, tt.startPort, tt.endPort, tt.top, tt.timeout, tt.concurrency)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFlags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
