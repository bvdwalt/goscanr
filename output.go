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

func printPortTable(w io.Writer, results []scanner.PortResult) {
	portW, stateW, serviceW := len("PORT"), len("STATE"), len("SERVICE")
	for _, r := range results {
		if width := len(r.Port) + 1 + len(r.Proto); width > portW {
			portW = width
		}
		if width := len(r.Service); width > serviceW {
			serviceW = width
		}
	}

	top := fmt.Sprintf("┌─%s─┬─%s─┬─%s─┐", strings.Repeat("─", portW), strings.Repeat("─", stateW), strings.Repeat("─", serviceW))
	mid := fmt.Sprintf("├─%s─┼─%s─┼─%s─┤", strings.Repeat("─", portW), strings.Repeat("─", stateW), strings.Repeat("─", serviceW))
	bot := fmt.Sprintf("└─%s─┴─%s─┴─%s─┘", strings.Repeat("─", portW), strings.Repeat("─", stateW), strings.Repeat("─", serviceW))

	fmt.Fprintln(w, top)
	fmt.Fprintf(w, "│ %s%-*s%s │ %s%-*s%s │ %s%-*s%s │\n",
		colorBold, portW, "PORT", colorReset,
		colorBold, stateW, "STATE", colorReset,
		colorBold, serviceW, "SERVICE", colorReset,
	)
	fmt.Fprintln(w, mid)

	for _, r := range results {
		port := r.Port + "/" + r.Proto
		fmt.Fprintf(w, "│ %s%-*s%s │ %s%-*s%s │ %-*s │\n",
			colorCyan, portW, port, colorReset,
			colorGreen, stateW, r.State, colorReset,
			serviceW, r.Service,
		)
	}

	fmt.Fprintln(w, bot)
}

func printPlainPorts(w io.Writer, ports []int) {
	for _, port := range ports {
		fmt.Fprintf(w, "%s%d%s is open\n", colorCyan, port, colorReset)
	}
}

func printHeader(w io.Writer, target string, ips []string, startPort, endPort int) {
	fmt.Fprintf(w, "%s%s%s (%s) — ports %s%d-%d%s\n",
		colorBold, target, colorReset,
		strings.Join(ips, ", "),
		colorYellow, startPort, endPort, colorReset,
	)
}
