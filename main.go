package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"time"

	"bvdwalt/goscanr/scanner"
)

func main() {
	target := flag.String("target", "", "Target host to scan")
	startPort := flag.Int("start", 1, "Start port")
	endPort := flag.Int("end", 1024, "End port")
	timeout := flag.Int("timeout", 300, "Timeout in milliseconds")
	concurrency := flag.Int("concurrency", 1000, "Maximum number of concurrent port scans")
	flag.Parse()

	if err := validateFlags(*target, *startPort, *endPort, *timeout, *concurrency); err != nil {
		fmt.Println("Flags parsing error:", err)
		os.Exit(1)
	}

	ips, err := net.LookupHost(*target)
	if err != nil {
		fmt.Printf("Failed to resolve %s: %v\n", *target, err)
		os.Exit(1)
	}

	fmt.Printf("Scanning %s (%v) from port %d to %d...\n", *target, ips, *startPort, *endPort)
	start := time.Now()

	found := scanner.Scan(ips, *startPort, *endPort, time.Duration(*timeout)*time.Millisecond, *concurrency)

	sort.Ints(found)
	for _, port := range found {
		fmt.Printf("Port %d is open\n", port)
	}

	fmt.Println("Duration:", time.Since(start))
	fmt.Println("Scan complete")
}
