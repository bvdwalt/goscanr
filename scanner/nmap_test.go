package scanner

import (
	"testing"
)

func TestBuildNmapArgs(t *testing.T) {
	args := buildNmapArgs([]string{"192.168.1.1", "192.168.1.2"}, []int{22, 80, 443})
	expected := []string{"-oX", "-", "-p", "22,80,443", "192.168.1.1", "192.168.1.2"}

	if len(args) != len(expected) {
		t.Fatalf("expected %d args, got %d: %v", len(expected), len(args), args)
	}
	for i, a := range args {
		if a != expected[i] {
			t.Errorf("arg[%d]: expected %q, got %q", i, expected[i], a)
		}
	}
}

func TestParseNmapXML(t *testing.T) {
	xml := []byte(`<?xml version="1.0"?>
<nmaprun>
  <host>
    <address addr="192.168.1.1" addrtype="ipv4"/>
    <ports>
      <port protocol="tcp" portid="22">
        <state state="open"/>
        <service name="ssh"/>
      </port>
      <port protocol="tcp" portid="80">
        <state state="open"/>
        <service name="http"/>
      </port>
    </ports>
  </host>
</nmaprun>`)

	results, err := parseNmapXML(xml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	expected := []PortResult{
		{IP: "192.168.1.1", Port: "22", Proto: "tcp", State: "open", Service: "ssh"},
		{IP: "192.168.1.1", Port: "80", Proto: "tcp", State: "open", Service: "http"},
	}
	for i, r := range results {
		if r != expected[i] {
			t.Errorf("result[%d]: expected %+v, got %+v", i, expected[i], r)
		}
	}
}

func TestParseNmapXML_FiltersNonOpen(t *testing.T) {
	xml := []byte(`<?xml version="1.0"?>
<nmaprun>
  <host>
    <address addr="192.168.1.1" addrtype="ipv4"/>
    <ports>
      <port protocol="tcp" portid="22">
        <state state="open"/>
        <service name="ssh"/>
      </port>
      <port protocol="tcp" portid="80">
        <state state="closed"/>
        <service name="http"/>
      </port>
      <port protocol="tcp" portid="443">
        <state state="filtered"/>
        <service name="https"/>
      </port>
    </ports>
  </host>
</nmaprun>`)

	results, err := parseNmapXML(xml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Port != "22" {
		t.Errorf("expected only open port 22, got %v", results)
	}
}

func TestParseNmapXML_Empty(t *testing.T) {
	xml := []byte(`<?xml version="1.0"?><nmaprun></nmaprun>`)

	results, err := parseNmapXML(xml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %v", results)
	}
}

func TestParseNmapXML_Invalid(t *testing.T) {
	_, err := parseNmapXML([]byte(`not xml`))
	if err == nil {
		t.Error("expected error for invalid XML, got nil")
	}
}
