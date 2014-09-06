package main

import (
	"flag"
	"log"
	"net"

	"github.com/davecheney/partyline"
)

var (
	endpoint string

	tunnel = make(map[partyline.MAC]*net.UDPAddr)
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
		f := partyline.Frame(frame.buf)
		src := partyline.SourceMAC(f)
		tunnel[src] = frame.src
		log.Printf("%v: %v", frame.src, f)
		for k, a := range tunnel {
			if k == src {
				// skip
				continue
			}
			log.Printf("relay: %v -> %v", f, a)
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
