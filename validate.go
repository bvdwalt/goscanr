package main

import "errors"

func validateFlags(target string, startPort, endPort, timeout, concurrency int) error {
	if target == "" {
		return errors.New("target is required")
	}
	if startPort < 1 || startPort > 65535 {
		return errors.New("start port must be between 1 and 65535")
	}
	if endPort < 1 || endPort > 65535 {
		return errors.New("end port must be between 1 and 65535")
	}
	if startPort > endPort {
		return errors.New("start port must be less than or equal to end port")
	}
	if timeout < 1 {
		return errors.New("timeout must be at least 1ms")
	}
	if concurrency < 1 {
		return errors.New("concurrency must be at least 1")
	}
	return nil
}
