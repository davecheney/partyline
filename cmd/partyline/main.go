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
	vin      int
)

func init() {
	flag.StringVar(&endpoint, "e", "partyline.dfc.io:9000", "partyline endpoint")
	flag.IntVar(&vin, "v", 1<<12, "vxlan identifier")

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
		buf = partyline.Encapsulate(buf, vin)
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
		_, buf = partyline.Deencapsulate(buf)
		log.Printf("tap.Write: % x", buf)
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
		buf := make([]byte, 1500)
		n, err := tap.Read(buf)
		check(err)
		buf = buf[:n]
		log.Printf("tap.Read: % x", buf)
		out <- buf[:n]
	}
}
