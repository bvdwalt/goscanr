package scanner

import "testing"

func TestTopPorts_ReturnsCorrectCount(t *testing.T) {
	for _, n := range []int{10, 100, 500} {
		ports, err := TopPorts(n)
		if err != nil {
			t.Fatalf("TopPorts(%d) unexpected error: %v", n, err)
		}
		if len(ports) != n {
			t.Errorf("TopPorts(%d) returned %d ports", n, len(ports))
		}
	}
}

func TestTopPorts_TooLow(t *testing.T) {
	_, err := TopPorts(0)
	if err == nil {
		t.Error("expected error for n=0")
	}
}

func TestTopPorts_TooHigh(t *testing.T) {
	_, err := TopPorts(len(topPorts) + 1)
	if err == nil {
		t.Error("expected error for n > list length")
	}
}

func TestTopPorts_ContainsCommonPorts(t *testing.T) {
	ports, err := TopPorts(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	common := map[int]bool{80: false, 443: false, 22: false}
	for _, p := range ports {
		if _, ok := common[p]; ok {
			common[p] = true
		}
	}
	for port, found := range common {
		if !found {
			t.Errorf("expected port %d in top 10", port)
		}
	}
}

func TestPortRange(t *testing.T) {
	ports := PortRange(80, 85)
	expected := []int{80, 81, 82, 83, 84, 85}
	if len(ports) != len(expected) {
		t.Fatalf("expected %d ports, got %d", len(expected), len(ports))
	}
	for i, p := range ports {
		if p != expected[i] {
			t.Errorf("port[%d]: expected %d, got %d", i, expected[i], p)
		}
	}
}
