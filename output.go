package main

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"bvdwalt/goscanr/scanner"
)

const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
)

func bannerFor(ip string, port int, scanResults []scanner.ScanResult) string {
	for _, r := range scanResults {
		if r.IP == ip && r.Port == port && r.Banner != "" {
			return r.Banner
		}
	}
	return ""
}

func printPortTable(w io.Writer, results []scanner.PortResult, scanResults []scanner.ScanResult) {
	ipW, portW, stateW, serviceW, bannerW := len("IP"), len("PORT"), len("STATE"), len("SERVICE"), len("BANNER")
	for _, r := range results {
		if width := len(r.IP); width > ipW {
			ipW = width
		}
		if width := len(r.Port) + 1 + len(r.Proto); width > portW {
			portW = width
		}
		if width := len(r.State); width > stateW {
			stateW = width
		}
		if width := len(r.Service); width > serviceW {
			serviceW = width
		}
		port := 0
		fmt.Sscanf(r.Port, "%d", &port)
		if width := len(bannerFor(r.IP, port, scanResults)); width > bannerW {
			bannerW = width
		}
	}

	sep := func(l, m, r string) string {
		return fmt.Sprintf("%s─%s─%s─%s─%s─%s─%s─%s─%s─%s─%s",
			l, strings.Repeat("─", ipW),
			m, strings.Repeat("─", portW),
			m, strings.Repeat("─", stateW),
			m, strings.Repeat("─", serviceW),
			m, strings.Repeat("─", bannerW), r)
	}

	fmt.Fprintln(w, sep("┌", "┬", "┐"))
	fmt.Fprintf(w, "│ %s%-*s%s │ %s%-*s%s │ %s%-*s%s │ %s%-*s%s │ %s%-*s%s │\n",
		colorBold, ipW, "IP", colorReset,
		colorBold, portW, "PORT", colorReset,
		colorBold, stateW, "STATE", colorReset,
		colorBold, serviceW, "SERVICE", colorReset,
		colorBold, bannerW, "BANNER", colorReset,
	)
	fmt.Fprintln(w, sep("├", "┼", "┤"))

	for _, r := range results {
		port := 0
		fmt.Sscanf(r.Port, "%d", &port)
		banner := bannerFor(r.IP, port, scanResults)
		fmt.Fprintf(w, "│ %-*s │ %s%-*s%s │ %s%-*s%s │ %-*s │ %-*s │\n",
			ipW, r.IP,
			colorCyan, portW, r.Port+"/"+r.Proto, colorReset,
			colorGreen, stateW, r.State, colorReset,
			serviceW, r.Service,
			bannerW, banner,
		)
	}

	fmt.Fprintln(w, sep("└", "┴", "┘"))
}

func printResults(w io.Writer, target string, found []scanner.ScanResult) {
	sort.Slice(found, func(i, j int) bool {
		if found[i].IP != found[j].IP {
			return found[i].IP < found[j].IP
		}
		return found[i].Port < found[j].Port
	})

	ports := make([]int, len(found))
	for i, r := range found {
		ports[i] = r.Port
	}

	var portResults []scanner.PortResult
	if scanner.NmapAvailable() && len(found) > 0 {
		uniqueIPs := uniqueIPs(found)
		var err error
		portResults, err = scanner.RunNmap(uniqueIPs, ports)
		if err != nil {
			fmt.Fprintf(w, "nmap error: %v\n", err)
		}
	} else {
		for _, r := range found {
			portResults = append(portResults, scanner.PortResult{
				IP:    r.IP,
				Port:  fmt.Sprintf("%d", r.Port),
				Proto: r.Proto,
				State: "open",
			})
		}
	}
	printPortTable(w, portResults, found)
}

func uniqueIPs(results []scanner.ScanResult) []string {
	seen := make(map[string]bool)
	var ips []string
	for _, r := range results {
		if !seen[r.IP] {
			seen[r.IP] = true
			ips = append(ips, r.IP)
		}
	}
	return ips
}

func printHeader(w io.Writer, target string, ips []string, top, startPort, endPort int) {
	var portDesc string
	if top > 0 {
		portDesc = fmt.Sprintf("top %d", top)
	} else {
		portDesc = fmt.Sprintf("%d-%d", startPort, endPort)
	}

	var ipDesc string
	if len(ips) > 4 {
		ipDesc = fmt.Sprintf("%s ... %s (%d hosts)", ips[0], ips[len(ips)-1], len(ips))
	} else {
		ipDesc = strings.Join(ips, ", ")
	}

	fmt.Fprintf(w, "%s%s%s (%s) — ports %s%s%s\n",
		colorBold, target, colorReset,
		ipDesc,
		colorYellow, portDesc, colorReset,
	)
}
