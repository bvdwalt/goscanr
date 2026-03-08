package main

import (
	"fmt"
	"io"
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

func bannerFor(port int, scanResults []scanner.ScanResult) string {
	for _, r := range scanResults {
		if r.Port == port && r.Banner != "" {
			return r.Banner
		}
	}
	return ""
}

func printPortTable(w io.Writer, results []scanner.PortResult, scanResults []scanner.ScanResult) {
	portW, stateW, serviceW, bannerW := len("PORT"), len("STATE"), len("SERVICE"), len("BANNER")
	for _, r := range results {
		if width := len(r.Port) + 1 + len(r.Proto); width > portW {
			portW = width
		}
		if width := len(r.Service); width > serviceW {
			serviceW = width
		}
		port := 0
		fmt.Sscanf(r.Port, "%d", &port)
		if width := len(bannerFor(port, scanResults)); width > bannerW {
			bannerW = width
		}
	}

	sep := func(l, m, r string) string {
		return fmt.Sprintf("%s─%s─%s─%s─%s─%s─%s─%s─%s",
			l, strings.Repeat("─", portW),
			m, strings.Repeat("─", stateW),
			m, strings.Repeat("─", serviceW),
			m, strings.Repeat("─", bannerW), r)
	}

	fmt.Fprintln(w, sep("┌", "┬", "┐"))
	fmt.Fprintf(w, "│ %s%-*s%s │ %s%-*s%s │ %s%-*s%s │ %s%-*s%s │\n",
		colorBold, portW, "PORT", colorReset,
		colorBold, stateW, "STATE", colorReset,
		colorBold, serviceW, "SERVICE", colorReset,
		colorBold, bannerW, "BANNER", colorReset,
	)
	fmt.Fprintln(w, sep("├", "┼", "┤"))

	for _, r := range results {
		port := 0
		fmt.Sscanf(r.Port, "%d", &port)
		banner := bannerFor(port, scanResults)
		fmt.Fprintf(w, "│ %s%-*s%s │ %s%-*s%s │ %-*s │ %-*s │\n",
			colorCyan, portW, r.Port+"/"+r.Proto, colorReset,
			colorGreen, stateW, r.State, colorReset,
			serviceW, r.Service,
			bannerW, banner,
		)
	}

	fmt.Fprintln(w, sep("└", "┴", "┘"))
}


func printHeader(w io.Writer, target string, ips []string, startPort, endPort int) {
	fmt.Fprintf(w, "%s%s%s (%s) — ports %s%d-%d%s\n",
		colorBold, target, colorReset,
		strings.Join(ips, ", "),
		colorYellow, startPort, endPort, colorReset,
	)
}
