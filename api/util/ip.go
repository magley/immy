package util

import (
	"encoding/binary"
	"net"
)

// Parse IPv4 address (d.d.d.d).
func IPv4toUint64(s string) uint64 {
	ip := net.ParseIP(s)
	if ip == nil {
		return 0
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return 0
	}

	return uint64(binary.BigEndian.Uint32(ipv4))
}