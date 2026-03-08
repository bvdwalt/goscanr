package scanner

import (
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ScanResult struct {
	Port   int
	Proto  string
	Banner string
}

func Scan(ips []string, startPort, endPort int, timeout time.Duration, concurrency int) []ScanResult {
	results := make(chan ScanResult)
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for _, ip := range ips {
		for port := startPort; port <= endPort; port++ {
			wg.Add(1)
			sem <- struct{}{}
			go func(ip string, p int) {
				defer func() { <-sem }()
				scanPort(ip, p, timeout, &wg, results)
			}(ip, port)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var found []ScanResult
	for r := range results {
		found = append(found, r)
	}
	return found
}

func scanPort(ip string, port int, timeout time.Duration, wg *sync.WaitGroup, results chan ScanResult) {
	defer wg.Done()

	address := net.JoinHostPort(ip, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return
	}
	defer conn.Close()

	banner := grabBanner(conn)
	results <- ScanResult{Port: port, Proto: "tcp", Banner: banner}
}

const bannerTimeout = 100 * time.Millisecond

func grabBanner(conn net.Conn) string {
	conn.SetReadDeadline(time.Now().Add(bannerTimeout))
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	return strings.TrimSpace(string(buf[:n]))
}
