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

	var buf strings.Builder
	printPortTable(&buf, results)
	out := buf.String()

	for _, want := range []string{"PORT", "STATE", "SERVICE", "22/tcp", "ssh", "443/tcp", "https"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestPrintPortTable_ColumnsAdjustToContent(t *testing.T) {
	results := []scanner.PortResult{
		{Port: "22", Proto: "tcp", State: "open", Service: "a-very-long-service-name"},
	}

	var buf strings.Builder
	printPortTable(&buf, results)

	for _, line := range strings.Split(buf.String(), "\n") {
		if strings.HasPrefix(line, "│") && strings.Contains(line, "a-very-long-service-name") {
			if !strings.HasSuffix(strings.TrimRight(line, " "), "│") {
				t.Error("table row is not properly closed")
			}
		}
	}
}

func TestPrintPlainPorts(t *testing.T) {
	var buf strings.Builder
	printPlainPorts(&buf, []int{22, 80, 443})
	out := buf.String()

	for _, want := range []string{"22", "80", "443"} {
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
