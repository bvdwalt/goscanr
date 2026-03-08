package scanner

import (
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	minConcurrency     = 50
	maxConcurrency     = 5000
	timeoutThresholdHi = 0.10 // reduce concurrency if >10% of batch timed out
	timeoutThresholdLo = 0.05 // increase concurrency if <5% of batch timed out
)

type ScanResult struct {
	Port   int
	Proto  string
	Banner string
}

type portTarget struct {
	ip   string
	port int
}

type scanOutcome struct {
	result   *ScanResult
	timedOut bool
}

func Scan(ips []string, startPort, endPort int, timeout time.Duration, concurrency int) []ScanResult {
	var targets []portTarget
	for _, ip := range ips {
		for port := startPort; port <= endPort; port++ {
			targets = append(targets, portTarget{ip, port})
		}
	}

	var found []ScanResult
	for i := 0; i < len(targets); {
		end := i + concurrency
		if end > len(targets) {
			end = len(targets)
		}
		batch := targets[i:end]
		i = end

		outcomes := make(chan scanOutcome, len(batch))
		var wg sync.WaitGroup
		for _, t := range batch {
			wg.Add(1)
			go func(ip string, port int) {
				defer wg.Done()
				scanPort(ip, port, timeout, outcomes)
			}(t.ip, t.port)
		}
		wg.Wait()
		close(outcomes)

		var timeouts int
		for o := range outcomes {
			if o.timedOut {
				timeouts++
			} else if o.result != nil {
				found = append(found, *o.result)
			}
		}

		concurrency = adjustConcurrency(concurrency, timeouts, len(batch))
	}

	return found
}

func adjustConcurrency(concurrency, timeouts, total int) int {
	rate := float64(timeouts) / float64(total)
	if rate > timeoutThresholdHi {
		concurrency = max(minConcurrency, concurrency*9/10) // reduce by 10%
	} else if rate < timeoutThresholdLo {
		concurrency = min(maxConcurrency, concurrency*21/20) // increase by 5%
	}
	return concurrency
}

func scanPort(ip string, port int, timeout time.Duration, outcomes chan scanOutcome) {
	address := net.JoinHostPort(ip, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		timedOut := false
		if netErr, ok := err.(net.Error); ok {
			timedOut = netErr.Timeout()
		}
		outcomes <- scanOutcome{timedOut: timedOut}
		return
	}
	defer conn.Close()

	banner := grabBanner(conn)
	outcomes <- scanOutcome{result: &ScanResult{Port: port, Proto: "tcp", Banner: banner}}
}

const bannerTimeout = 100 * time.Millisecond

func grabBanner(conn net.Conn) string {
	conn.SetReadDeadline(time.Now().Add(bannerTimeout))
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	return strings.TrimSpace(string(buf[:n]))
}
