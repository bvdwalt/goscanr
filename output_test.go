package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"bvdwalt/goscanr/scanner"
)

func TestPrintPortTable(t *testing.T) {
	results := []scanner.PortResult{
		{Port: "22", Proto: "tcp", State: "open", Service: "ssh"},
		{Port: "443", Proto: "tcp", State: "open", Service: "https"},
	}
	scanResults := []scanner.ScanResult{
		{Port: 22, Banner: "SSH-2.0-OpenSSH_9.3"},
	}

	var buf strings.Builder
	printPortTable(&buf, results, scanResults)
	out := buf.String()

	for _, want := range []string{"PORT", "STATE", "SERVICE", "BANNER", "22/tcp", "ssh", "443/tcp", "https", "SSH-2.0-OpenSSH_9.3"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestPrintPortTable_EmptyBanner(t *testing.T) {
	results := []scanner.PortResult{
		{Port: "80", Proto: "tcp", State: "open", Service: "http"},
	}

	var buf strings.Builder
	printPortTable(&buf, results, nil)
	out := buf.String()

	if !strings.Contains(out, "80/tcp") {
		t.Error("expected output to contain '80/tcp'")
	}
}


func TestPrintHeader(t *testing.T) {
	var buf strings.Builder
	printHeader(&buf, "example.com", []string{"1.2.3.4"}, 0, 1, 1024)
	out := buf.String()

	for _, want := range []string{"example.com", "1.2.3.4", "1-1024"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestPrintHeader_Top(t *testing.T) {
	var buf strings.Builder
	printHeader(&buf, "example.com", []string{"1.2.3.4"}, 100, 0, 0)
	out := buf.String()

	if !strings.Contains(out, "top 100") {
		t.Errorf("expected output to contain 'top 100', got: %s", out)
	}
}

func TestPrintHeader_ManyIPs(t *testing.T) {
	ips := []string{"1.1.1.1", "1.1.1.2", "1.1.1.3", "1.1.1.4", "1.1.1.5"}
	var buf strings.Builder
	printHeader(&buf, "10.0.0.0/24", ips, 0, 1, 1024)
	out := buf.String()

	for _, want := range []string{"1.1.1.1", "1.1.1.5", "5 hosts"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got: %s", want, out)
		}
	}
}

func TestUniqueIPs(t *testing.T) {
	results := []scanner.ScanResult{
		{IP: "1.1.1.1", Port: 80},
		{IP: "1.1.1.1", Port: 443},
		{IP: "1.1.1.2", Port: 80},
	}
	ips := uniqueIPs(results)
	if len(ips) != 2 {
		t.Errorf("expected 2 unique IPs, got %d", len(ips))
	}
}

func TestPrintResultsJSON(t *testing.T) {
	portResults := []scanner.PortResult{
		{IP: "1.2.3.4", Port: "80", Proto: "tcp", State: "open", Service: "http"},
	}
	scanResults := []scanner.ScanResult{
		{IP: "1.2.3.4", Port: 80, Proto: "tcp", Banner: "nginx"},
	}
	var buf strings.Builder
	printResultsJSON(&buf, "example.com", portResults, scanResults, 200*time.Millisecond)

	var out jsonOutput
	if err := json.Unmarshal([]byte(buf.String()), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}
	if out.Target != "example.com" {
		t.Errorf("expected target 'example.com', got %q", out.Target)
	}
	if len(out.Results) == 0 {
		t.Fatal("expected at least one result")
	}
	if out.Results[0].Port != "80" {
		t.Errorf("expected port '80', got %q", out.Results[0].Port)
	}
	if out.Results[0].Banner != "nginx" {
		t.Errorf("expected banner 'nginx', got %q", out.Results[0].Banner)
	}
}

func TestPrintResults_DefaultIsText(t *testing.T) {
	var buf strings.Builder
	// empty found means nmap is skipped, portResults are empty too — just verify no JSON output
	printResults(&buf, "example.com", nil, "text", 100*time.Millisecond)
	out := buf.String()
	// text format renders a table (box drawing characters), not JSON
	if strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Error("expected text format, got JSON-like output")
	}
}

func TestAnsiStripper(t *testing.T) {
	var buf strings.Builder
	w := ansiStripper{&buf}
	input := "\033[1mhello\033[0m world"
	w.Write([]byte(input))
	if got := buf.String(); got != "hello world" {
		t.Errorf("expected ANSI stripped output, got %q", got)
	}
}
