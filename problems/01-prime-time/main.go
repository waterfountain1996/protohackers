package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"time"
)

const IsPrimeMethod = "isPrime"

var MalformedResponse = []byte("418 I'm a teapot")

var MalformedResponseError = fmt.Errorf("Malformed response error")

var ReadTimeout = time.Duration(5 * 1e9)

type Request struct {
	Method string  `json:"method"`
	Number float64 `json:"number"`
}

func NewRequestFromBytes(data []byte) (*Request, error) {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}

	switch v["method"].(type) {
	case string:
	default:
		return nil, MalformedResponseError
	}

	switch v["number"].(type) {
	case float64:
	default:
		return nil, MalformedResponseError
	}

	return &Request{
		Method: v["method"].(string),
		Number: v["number"].(float64),
	}, nil
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func NewResponse(prime bool) *Response {
	return &Response{
		Method: "isPrime",
		Prime:  prime,
	}
}

func isPrime(number int64) bool {
	if number <= 1 {
		return false
	}

	limit := int64(math.Sqrt(float64(number)))

	for i := int64(2); i <= limit; i++ {
		if number%i == 0 {
			return false
		}
	}

	return true
}

func connHandler(conn net.Conn) {
	defer conn.Close()

	s := bufio.NewScanner(conn)
	enc := json.NewEncoder(conn)

	conn.SetReadDeadline(time.Now().Add(ReadTimeout))

	for s.Scan() {
		req, err := NewRequestFromBytes(s.Bytes())
		if err != nil || req.Method != IsPrimeMethod {
			log.Printf("Malformed connection from %s\n", conn.RemoteAddr())
			conn.Write(MalformedResponse)
			return
		}

		var prime bool

		if req.Number != math.Trunc(req.Number) {
			prime = false
		} else {
			prime = isPrime(int64(req.Number))
		}

		res := NewResponse(prime)
		enc.Encode(res)
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
