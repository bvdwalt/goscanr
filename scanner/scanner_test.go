package scanner

import (
	"net"
	"sort"
	"testing"
	"time"
)

func TestScan_FindsOpenPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port

	found := Scan([]string{"127.0.0.1"}, port, port, 500*time.Millisecond, 1)

	if len(found) != 1 || found[0] != port {
		t.Errorf("expected [%d], got %v", port, found)
	}
}

func TestScan_NoOpenPorts(t *testing.T) {
	found := Scan([]string{"127.0.0.1"}, 1, 1, 50*time.Millisecond, 1)

	if len(found) != 0 {
		t.Errorf("expected no open ports, got %v", found)
	}
}

func TestScan_ReturnsMultipleOpenPorts(t *testing.T) {
	var listeners []net.Listener
	var ports []int
	for range 3 {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("failed to start listener: %v", err)
		}
		defer ln.Close()
		listeners = append(listeners, ln)
		ports = append(ports, ln.Addr().(*net.TCPAddr).Port)
	}

	sort.Ints(ports)
	found := Scan([]string{"127.0.0.1"}, ports[0], ports[len(ports)-1], 500*time.Millisecond, 10)
	sort.Ints(found)

	for _, p := range ports {
		if !contains(found, p) {
			t.Errorf("expected port %d in results %v", p, found)
		}
	}
}

func contains(s []int, v int) bool {
	for _, n := range s {
		if n == v {
			return true
		}
	}
	return false
}
