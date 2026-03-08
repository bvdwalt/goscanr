package main

import (
	"fmt"
	"os"
)

func validateFlags(target string, startPort, endPort, timeout, concurrency int) {
	if target == "" {
		fmt.Println("Usage: main -target <host> -start <startPort> -end <endPort>")
		os.Exit(1)
	}
	if startPort < 1 || startPort > 65535 {
		fmt.Println("Error: start port must be between 1 and 65535")
		os.Exit(1)
	}
	if endPort < 1 || endPort > 65535 {
		fmt.Println("Error: end port must be between 1 and 65535")
		os.Exit(1)
	}
	if startPort > endPort {
		fmt.Println("Error: start port must be less than or equal to end port")
		os.Exit(1)
	}
	if timeout < 1 {
		fmt.Println("Error: timeout must be at least 1ms")
		os.Exit(1)
	}
	if concurrency < 1 {
		fmt.Println("Error: concurrency must be at least 1")
		os.Exit(1)
	}
}
