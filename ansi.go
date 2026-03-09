package main

import (
	"io"
	"regexp"
)

var ansiEscape = regexp.MustCompile(`\033\[[0-9;]*m`)

type ansiStripper struct {
	w io.Writer
}

func (a ansiStripper) Write(p []byte) (n int, err error) {
	stripped := ansiEscape.ReplaceAll(p, nil)
	_, err = a.w.Write(stripped)
	return len(p), err
}
