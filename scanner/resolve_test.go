package scanner

import "testing"

func TestResolveTarget_CIDR(t *testing.T) {
	ips, err := ResolveTarget("192.168.1.0/24")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ips) != 254 {
		t.Errorf("expected 254 IPs, got %d", len(ips))
	}
}

func TestResolveTarget_Hostname(t *testing.T) {
	ips, err := ResolveTarget("localhost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ips) == 0 {
		t.Error("expected at least one IP for localhost")
	}
}

func TestResolveTarget_InvalidCIDR(t *testing.T) {
	_, err := ResolveTarget("192.168.1.0/99")
	if err == nil {
		t.Error("expected error for invalid CIDR")
	}
}
