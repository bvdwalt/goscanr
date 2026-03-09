package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"bvdwalt/goscanr/scanner"
)

func resolvePorts(top, start, end int) ([]int, error) {
	if top > 0 {
		return scanner.TopPorts(top)
	}
	return scanner.PortRange(start, end), nil
}

func main() {
	target := flag.String("target", "", "Target host to scan")
	startPort := flag.Int("start", 1, "Start port")
	endPort := flag.Int("end", 1024, "End port")
	top := flag.Int("top", 0, "Scan the top N most common ports (overrides -start and -end)")
	timeout := flag.Int("timeout", 300, "Timeout in milliseconds")
	concurrency := flag.Int("concurrency", 500, "Initial number of concurrent port scans (adapts automatically)")
	output := flag.String("output", "", "Save results to a file")
	format := flag.String("format", "text", "Output format: text or json")
	flag.Parse()

	if err := validateFlags(*target, *startPort, *endPort, *top, *timeout, *concurrency, *format); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if !scanner.NmapAvailable() {
		fmt.Println("Warning: nmap not found in PATH, skipping service detection")
	}

	ips, err := scanner.ResolveTarget(*target)
	if err != nil {
		fmt.Printf("Failed to resolve %s: %v\n", *target, err)
		os.Exit(1)
	}

	out := io.Writer(os.Stdout)
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			os.Exit(1)
		}
		defer f.Close()
		out = io.MultiWriter(out, ansiStripper{f})
	}

	ports, err := resolvePorts(*top, *startPort, *endPort)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if *format == "json" {
		printHeader(os.Stderr, *target, ips, *top, *startPort, *endPort)
	} else {
		printHeader(out, *target, ips, *top, *startPort, *endPort)
	}

	start := time.Now()
	found := scanner.Scan(ips, ports, time.Duration(*timeout)*time.Millisecond, *concurrency, func(done, total int) {
		fmt.Fprintf(os.Stderr, "\rScanning... %d/%d (%.0f%%)", done, total, float64(done)/float64(total)*100)
		if done == total {
			fmt.Fprintln(os.Stderr)
		}
	})
	duration := time.Since(start)
	printResults(out, *target, found, *format, duration)

	if *format == "json" {
		fmt.Fprintf(os.Stderr, "Done in %s\n", duration.Round(time.Millisecond))
	} else {
		fmt.Fprintf(out, "Done in %s\n", duration.Round(time.Millisecond))
	}
}
