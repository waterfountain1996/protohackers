package main

import (
	"bufio"
	"log"
	"net"
)

func readSpec(br *bufio.Reader) ([]byte, error) {
	var (
		spec      = make([]byte, 81)
		idx       = 0
		expectArg = false
	)

	for {
		b, err := br.ReadByte()
		if err != nil {
			return nil, err
		}

		spec[idx] = b

		if !expectArg {
			switch b {
			case 0x00:
				return spec[:idx+1], nil
			case 0x02:
				fallthrough
			case 0x04:
				expectArg = true
			}
		} else {
			expectArg = false
		}

		idx++
	}
}

func handleConn(conn net.Conn) {
	defer func() {
		log.Printf("%s disconnected\n", conn.RemoteAddr())
		conn.Close()
	}()

	r := bufio.NewReader(conn)

	spec, err := readSpec(r)
	if err != nil {
		log.Printf("Failed to read cipher spec from %s: %s\n", conn.RemoteAddr(), err)
		return
	}

	if isNoop(spec) {
		log.Printf("%s sent a no-op spec %v\n", conn.RemoteAddr(), spec)
		return
	}

	scanner := bufio.NewScanner(NewObfuscatedReader(r, spec))
	w := NewObfuscatedWriter(conn, spec)

	for scanner.Scan() {
		request := scanner.Text()
		response := FindToy(request)

		if _, err := w.Write(append([]byte(response), '\n')); err != nil {
			log.Printf("Write error: %s\n", err)
			break
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":10000")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("New connection from %s\n", conn.RemoteAddr())

		go handleConn(conn)
	}
}
