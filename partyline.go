// Package partyline handles encapsulation and deencapsulation of ethernet frames
// in VXLAN frames.
package partyline

func Encapsulate(buf []byte, vin int) []byte {
	vin &= (1<<24 - 1)
	header := [8]byte{
		0: 0x08,
		4: uint8(vin >> 16),
		5: uint8(vin >> 8),
		6: uint8(vin),
	}
	return append(header[:], buf...)
}

func Deencapsulate(buf []byte) (int, []byte) {
	header := buf[:8]
	vin := int(header[4])<<16 | int(header[5])<<8 | int(header[6])
	return vin, buf[8:]
}
