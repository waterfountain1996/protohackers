package main

import (
	"io"
	"log"
	"net"
)

// TCP Echo Service from RFC 862
func connHandler(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 4096)
	if _, err := io.CopyBuffer(conn, conn, buffer); err != nil {
		log.Fatal(err)
	}

	log.Printf("%s disconnected\n", conn.RemoteAddr())
}

func main() {
	ln, err := net.Listen("tcp", ":10000")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s\n", ":10000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("TCP connection from %s\n", conn.RemoteAddr())

		go connHandler(conn)
	}
}
