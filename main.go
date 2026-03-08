package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

func main() {
	target := flag.String("target", "", "Target host to scan")
	startPort := flag.Int("start", 1, "Start port")
	endPort := flag.Int("end", 1024, "End port")
	timeout := flag.Int("timeout", 300, "Timeout in milliseconds")
	concurrency := flag.Int("concurrency", 1000, "Maximum number of concurrent port scans")
	flag.Parse()

	if *target == "" {
		fmt.Println("Usage: main -target <host> -start <startPort> -end <endPort>")
		os.Exit(1)
	}

	ips, err := net.LookupHost(*target)
	if err != nil {
		fmt.Printf("Failed to resolve %s: %v\n", *target, err)
		os.Exit(1)
	}

	fmt.Printf("Scanning %s (%v) from port %d to %d...\n", *target, ips, *startPort, *endPort)
	start := time.Now()

	openPorts := make(chan int)
	sem := make(chan struct{}, *concurrency)
	var wg sync.WaitGroup
	for _, targetIP := range ips {
		for port := *startPort; port <= *endPort; port++ {
			wg.Add(1)
			sem <- struct{}{}
			go func(ip string, p int) {
				defer func() { <-sem }()
				scanPort(ip, p, *timeout, &wg, openPorts)
			}(targetIP, port)
		}
	}

	go func() {
		wg.Wait()
		close(openPorts)
	}()

	var found []int
	for port := range openPorts {
		found = append(found, port)
	}

	sort.Ints(found)
	for _, port := range found {
		fmt.Printf("Port %d is open\n", port)
	}

	fmt.Println("Duration:", time.Since(start))
	fmt.Println("Scan complete")
}

func scanPort(target string, port int, timeout int, wg *sync.WaitGroup, openPorts chan int) {
	defer wg.Done()

	address := net.JoinHostPort(target, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", address, time.Millisecond*time.Duration(timeout))
	if err != nil {
		return
	}
	defer conn.Close()
	openPorts <- port
}
