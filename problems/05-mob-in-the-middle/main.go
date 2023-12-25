package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"regexp"
	"sync"
)

const (
	// Address of the upstream chat server
	upstreamAddress string = "chat.protohackers.com:16963"

	// Tony's boguscoin address
	addressReplace string = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
)

var (
	// Regexp that matches boguscoin addresses
	addressRegexp = regexp.MustCompile(`(?:^|\s)([7]{1}[[:alnum:]]{25,34})(?:$|\s)`)
)

func proxyRewrite(src io.Reader, dst io.Writer, rewriter func([]byte) []byte) {
	r := bufio.NewReader(src)
	for {
		data, err := r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %s\n", err)
			}
			break
		}

		out := rewriter(data)
		if _, err := dst.Write(append(out, '\n')); err != nil {
			log.Printf("Write error: %s\n", err)
			break
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	remote, err := net.Dial("tcp", upstreamAddress)
	if err != nil {
		log.Fatalf("Failed to connect to upstream server: %s\n", err)
	}

	rewriter := func(data []byte) []byte {
		s := bufio.NewScanner(bytes.NewReader(data))
		s.Split(bufio.ScanWords)

		words := [][]byte{}
		for s.Scan() {
			word := addressRegexp.ReplaceAll(s.Bytes(), []byte(addressReplace))
			words = append(words, word)
		}

		return bytes.Join(words, []byte(" "))
	}

	var wg sync.WaitGroup

	bg := func(src, dst net.Conn) {
		defer wg.Done()
		defer dst.Close()
		proxyRewrite(src, dst, rewriter)
	}

	wg.Add(2)
	go bg(remote, conn)
	go bg(conn, remote)

	wg.Wait()

	log.Printf("%s disconnected\n", conn.RemoteAddr())
}

func main() {
	ln, err := net.Listen("tcp", ":10000")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Printf("Listening on %s...\n", ":10000")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("New connection from %s\n", conn.RemoteAddr())

		go handleConn(conn)
	}
}
