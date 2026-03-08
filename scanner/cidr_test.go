package scanner

import (
	"testing"
)

func TestExpandCIDR_Slash24(t *testing.T) {
	ips, err := ExpandCIDR("192.168.1.0/24")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// /24 = 256 addresses, minus network and broadcast = 254
	if len(ips) != 254 {
		t.Errorf("expected 254 IPs, got %d", len(ips))
	}
	if ips[0] != "192.168.1.1" {
		t.Errorf("expected first IP to be 192.168.1.1, got %s", ips[0])
	}
	if ips[len(ips)-1] != "192.168.1.254" {
		t.Errorf("expected last IP to be 192.168.1.254, got %s", ips[len(ips)-1])
	}
}

func TestExpandCIDR_Slash31(t *testing.T) {
	ips, err := ExpandCIDR("10.0.0.0/31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// /31 has only 2 addresses, no network/broadcast stripped
	if len(ips) != 2 {
		t.Errorf("expected 2 IPs, got %d", len(ips))
	}
}

func TestExpandCIDR_Slash32(t *testing.T) {
	ips, err := ExpandCIDR("10.0.0.1/32")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ips) != 1 {
		t.Errorf("expected 1 IP, got %d", len(ips))
	}
	if ips[0] != "10.0.0.1" {
		t.Errorf("expected 10.0.0.1, got %s", ips[0])
	}
}

func TestExpandCIDR_Invalid(t *testing.T) {
	_, err := ExpandCIDR("not-a-cidr")
	if err == nil {
		t.Error("expected error for invalid CIDR, got nil")
	}
}
