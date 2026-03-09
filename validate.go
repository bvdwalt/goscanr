package main

import (
	"errors"
	"fmt"

	"bvdwalt/goscanr/scanner"
)

func validateFlags(target string, startPort, endPort, top, timeout, concurrency int, format string) error {
	if target == "" {
		return errors.New("target is required")
	}
	if top > 0 {
		if _, err := scanner.TopPorts(top); err != nil {
			return fmt.Errorf("invalid -top value: %w", err)
		}
	} else {
		if startPort < 1 || startPort > 65535 {
			return errors.New("start port must be between 1 and 65535")
		}
		if endPort < 1 || endPort > 65535 {
			return errors.New("end port must be between 1 and 65535")
		}
		if startPort > endPort {
			return errors.New("start port must be less than or equal to end port")
		}
	}
	if timeout < 1 {
		return errors.New("timeout must be at least 1ms")
	}
	if concurrency < 1 {
		return errors.New("concurrency must be at least 1")
	}
	if format != "text" && format != "json" {
		return fmt.Errorf("invalid format %q: must be text or json", format)
	}
	return nil
}
