package partyline

import (
	"fmt"
)

// MAC represents an IEEE 802 MAC address.
type MAC [6]uint8

func (m MAC) String() string {
	return fmt.Sprintf("%x:%x:%x:%x:%x:%x", m[0], m[1], m[2], m[3], m[4], m[5])
}

// SourceMAC returns the source MAC address of the frame.
func SourceMAC(b []byte) MAC {
	var m MAC
	copy(m[:], b[6:12])
	return m
}

// DestMAC returns the destination MAC address of the frame.
func DestMAC(b []byte) MAC {
	var m MAC
	copy(m[:], b[:6])
	return m
}

// A frame represents the contents of a complete ethernet frame.
type Frame []byte

func (f Frame) String() string {
	return fmt.Sprintf("%v -> %v length: %d", SourceMAC(f), DestMAC(f), len(f))
}
