package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"bvdwalt/goscanr/scanner"
)

func main() {
	target := flag.String("target", "", "Target host to scan")
	startPort := flag.Int("start", 1, "Start port")
	endPort := flag.Int("end", 1024, "End port")
	timeout := flag.Int("timeout", 300, "Timeout in milliseconds")
	concurrency := flag.Int("concurrency", 500, "Initial number of concurrent port scans (adapts automatically)")
	flag.Parse()

	if err := validateFlags(*target, *startPort, *endPort, *timeout, *concurrency); err != nil {
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

	printHeader(os.Stdout, *target, ips, *startPort, *endPort)
	start := time.Now()

	found := scanner.Scan(ips, *startPort, *endPort, time.Duration(*timeout)*time.Millisecond, *concurrency)
	printResults(os.Stdout, *target, found)

	fmt.Printf("Done in %s\n", time.Since(start).Round(time.Millisecond))
}
