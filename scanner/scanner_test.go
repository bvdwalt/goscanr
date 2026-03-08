package scanner

import (
	"net"
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
	found := Scan([]string{"127.0.0.1"}, port, port, 500*time.Millisecond, 10)

	if len(found) != 1 || found[0].Port != port {
		t.Errorf("expected port %d, got %v", port, found)
	}
}

func TestScan_NoOpenPorts(t *testing.T) {
	found := Scan([]string{"127.0.0.1"}, 1, 1, 50*time.Millisecond, 10)

	if len(found) != 0 {
		t.Errorf("expected no open ports, got %v", found)
	}
}

func TestScan_ReturnsMultipleOpenPorts(t *testing.T) {
	var ports []int
	for range 3 {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("failed to start listener: %v", err)
		}
		defer ln.Close()
		ports = append(ports, ln.Addr().(*net.TCPAddr).Port)
	}

	min, max := ports[0], ports[0]
	for _, p := range ports {
		if p < min {
			min = p
		}
		if p > max {
			max = p
		}
	}

	found := Scan([]string{"127.0.0.1"}, min, max, 500*time.Millisecond, 10)

	if len(found) < 3 {
		t.Errorf("expected at least 3 open ports, got %v", found)
	}
}

func TestScan_GrabsBanner(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		conn.Write([]byte("Hello from test server"))
	}()

	port := ln.Addr().(*net.TCPAddr).Port
	found := Scan([]string{"127.0.0.1"}, port, port, 500*time.Millisecond, 10)

	if len(found) != 1 {
		t.Fatalf("expected 1 result, got %d", len(found))
	}
	if found[0].Banner != "Hello from test server" {
		t.Errorf("expected banner %q, got %q", "Hello from test server", found[0].Banner)
	}
}

func TestAdjustConcurrency_ReducesOnHighTimeouts(t *testing.T) {
	result := adjustConcurrency(1000, 150, 1000) // 15% timeout rate
	if result >= 1000 {
		t.Errorf("expected concurrency to decrease, got %d", result)
	}
}

func TestAdjustConcurrency_IncreasesOnLowTimeouts(t *testing.T) {
	result := adjustConcurrency(1000, 10, 1000) // 1% timeout rate
	if result <= 1000 {
		t.Errorf("expected concurrency to increase, got %d", result)
	}
}

func TestAdjustConcurrency_RespectsMinimum(t *testing.T) {
	result := adjustConcurrency(minConcurrency, 1000, 1000) // 100% timeouts
	if result < minConcurrency {
		t.Errorf("expected concurrency to stay >= %d, got %d", minConcurrency, result)
	}
}

func TestAdjustConcurrency_RespectsMaximum(t *testing.T) {
	result := adjustConcurrency(maxConcurrency, 0, 1000) // 0% timeouts
	if result > maxConcurrency {
		t.Errorf("expected concurrency to stay <= %d, got %d", maxConcurrency, result)
	}
}
