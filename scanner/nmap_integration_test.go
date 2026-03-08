package scanner

import (
	"net"
	"testing"
)

func TestNmapAvailable(t *testing.T) {
	// Just verify it returns without panicking — result depends on the system
	_ = NmapAvailable()
}

func TestRunNmap(t *testing.T) {
	if !NmapAvailable() {
		t.Skip("nmap not installed, skipping")
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	results, err := RunNmap("127.0.0.1", []int{port})
	if err != nil {
		t.Fatalf("RunNmap returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].State != "open" {
		t.Errorf("expected state 'open', got %q", results[0].State)
	}
}
