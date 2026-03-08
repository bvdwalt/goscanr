package scanner

import (
	"net"
	"strconv"
	"sync"
	"time"
)

func Scan(ips []string, startPort, endPort int, timeout time.Duration, concurrency int) []int {
	openPorts := make(chan int)
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for _, ip := range ips {
		for port := startPort; port <= endPort; port++ {
			wg.Add(1)
			sem <- struct{}{}
			go func(ip string, p int) {
				defer func() { <-sem }()
				scanPort(ip, p, timeout, &wg, openPorts)
			}(ip, port)
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
	return found
}

func scanPort(ip string, port int, timeout time.Duration, wg *sync.WaitGroup, openPorts chan int) {
	defer wg.Done()

	address := net.JoinHostPort(ip, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return
	}
	defer conn.Close()
	openPorts <- port
}
