package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	target := flag.String("target", "", "Target host to scan")
	startPort := flag.Int("start", 1, "Start port")
	endPort := flag.Int("end", 1024, "End port")
	timeout := flag.Int("timeout", 300, "Timeout in milliseconds")
	flag.Parse()

	if *target == "" {
		fmt.Println("Usage: main -target <host> -start <startPort> -end <endPort>")
		os.Exit(1)
	}

	fmt.Printf("Scanning %s from port %d to %d...\n", *target, *startPort, *endPort)

	var wg sync.WaitGroup
	for port := *startPort; port <= *endPort; port++ {
		wg.Add(1)
		go scanPort(*target, port, *timeout, &wg)
	}
	wg.Wait()
	fmt.Println("Scan complete.")
}

func scanPort(target string, port int, timeout int, wg *sync.WaitGroup) {
	defer wg.Done()

	address := fmt.Sprintf("%s:%d", target, port)
	conn, err := net.DialTimeout("tcp", address, time.Millisecond*time.Duration(timeout))
	if err != nil {
		return
	}
	defer conn.Close()
	fmt.Printf("Port %d is open\n", port)
}
