package scanner

import (
	"net"
	"strings"
)

func ResolveTarget(target string) ([]string, error) {
	if strings.Contains(target, "/") {
		return ExpandCIDR(target)
	}
	return net.LookupHost(target)
}
