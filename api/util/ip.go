package util

import (
	"encoding/binary"
	"net"
)

// Parse IPv4 address (d.d.d.d).
func IPv4toUint64(s string) uint64 {
	var p [net.IPv4len]byte
	for i := 0; i < net.IPv4len; i++ {
		if len(s) == 0 {
			// Missing octets.
			return 0
		}
		if i > 0 {
			if s[0] != '.' {
				return 0
			}
			s = s[1:]
		}
		n, c, ok := dtoi(s)
		if !ok || n > 0xFF {
			return 0
		}
		s = s[c:]
		p[i] = byte(n)
	}

	if len(s) != 0 {
		return 0
	}

	return binary.BigEndian.Uint64(p[:4])
}

// Decimal to integer.
// Returns number, characters consumed, success.
func dtoi(s string) (n int, i int, ok bool) {
	const big = 0xFFFFFF
	n = 0
	for i = 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n >= big {
			return big, i, false
		}
	}

	if i == 0 {
		return 0, 0, false
	}

	return n, i, true
}