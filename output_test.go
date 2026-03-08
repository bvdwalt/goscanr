package main

import (
	"strings"
	"testing"

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

func TestPrintPlainPorts(t *testing.T) {
	results := []scanner.ScanResult{
		{Port: 22, Banner: "SSH-2.0-OpenSSH_9.3"},
		{Port: 80, Banner: ""},
	}

	var buf strings.Builder
	printPlainPorts(&buf, results)
	out := buf.String()

	for _, want := range []string{"22", "80", "SSH-2.0-OpenSSH_9.3"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestPrintHeader(t *testing.T) {
	var buf strings.Builder
	printHeader(&buf, "example.com", []string{"1.2.3.4"}, 1, 1024)
	out := buf.String()

	for _, want := range []string{"example.com", "1.2.3.4", "1-1024"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}
