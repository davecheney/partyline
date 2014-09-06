package main

import (
	"flag"
	"io"
	"log"
	"net"

	"github.com/davecheney/partyline"
	"github.com/pkg/taptun"
)

var (
	endpoint string
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

func send(out <-chan []byte, conn *net.UDPConn) {
	defer conn.Close()
	for buf := range out {
		_, err := conn.Write(buf)
		check(err)
	}
}

func receive(conn *net.UDPConn, tap io.Writer) {
	var b [1600]byte
	for {
		n, err := conn.Read(b[:])
		check(err)
		buf := b[:n]
		log.Printf("tap.Write: %v", partyline.Frame(buf))
		_, err = tap.Write(buf)
	}
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", endpoint)
	check(err)

	conn, err := net.DialUDP("udp", nil, addr)
	check(err)
	log.Println("connected to partyline endpoint:", addr)

	tap, err := taptun.OpenTap()
	check(err)
	defer tap.Close()
	log.Println("local tunnel device:", tap)

	out := make(chan []byte, 16)

	go send(out, conn)

	go receive(conn, tap)

	for {
		buf := make([]byte, 1560)
		buf, err := taptun.ReadFrame(tap, buf)
		check(err)
		log.Printf("tap.Read: %s", partyline.Frame(buf))

		out <- buf
	}
}
