package main

import (
	"flag"
	"log"
	"net"

	"github.com/davecheney/partyline"
)

var (
	endpoint string

	tunnels = make(map[int]map[string]*net.UDPAddr)
)

func init() {
	flag.StringVar(&endpoint, "e", "partyline.dfc.io:9000", "partyline endpoint")

	flag.Parse()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handle(ch chan struct {
	buf []byte
	src *net.UDPAddr
}, conn *net.UDPConn) {
	for frame := range ch {
		vin, _ := partyline.Deencapsulate(frame.buf)

		tunnel := tunnels[vin]
		if tunnel == nil {
			tunnel = make(map[string]*net.UDPAddr)
			tunnels[vin] = tunnel
		}
		tunnel[frame.src.String()] = frame.src
		log.Printf("vin: %0x, src: %v", vin, frame.src)
		for k, a := range tunnel {
			if k == frame.src.String() {
				// skip
				continue
			}
			log.Printf("relay: %v -> %v", frame.src, a)
			_, err := conn.WriteToUDP(frame.buf, a)
			check(err)
		}
	}
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", endpoint)
	check(err)

	conn, err := net.ListenUDP("udp", addr)
	check(err)
	log.Println("listening on partyline endpoint:", addr)

	in := make(chan struct {
		buf []byte
		src *net.UDPAddr
	}, 16)

	go handle(in, conn)

	for {
		buf := make([]byte, 1600)
		n, src, err := conn.ReadFromUDP(buf)
		check(err)
		in <- struct {
			buf []byte
			src *net.UDPAddr
		}{buf[:n], src}
	}
}
