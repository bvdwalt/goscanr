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
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if !scanner.NmapAvailable() {
		fmt.Println("Warning: nmap not found in PATH, skipping service detection")
	}

	ips, err := net.LookupHost(*target)
	if err != nil {
		fmt.Printf("Failed to resolve %s: %v\n", *target, err)
		os.Exit(1)
	}

	printHeader(os.Stdout, *target, ips, *startPort, *endPort)
	start := time.Now()

	found := scanner.Scan(ips, *startPort, *endPort, time.Duration(*timeout)*time.Millisecond, *concurrency)
	sort.Slice(found, func(i, j int) bool { return found[i].Port < found[j].Port })

	ports := make([]int, len(found))
	for i, r := range found {
		ports[i] = r.Port
	}

	if scanner.NmapAvailable() && len(found) > 0 {
		results, err := scanner.RunNmap(*target, ports)
		if err != nil {
			fmt.Println("nmap error:", err)
		}
		printPortTable(os.Stdout, results, found)
	} else {
		printPlainPorts(os.Stdout, found)
	}

	fmt.Printf("Done in %s\n", time.Since(start).Round(time.Millisecond))
}
