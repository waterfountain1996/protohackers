package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/waterfountain1996/protohackers/datastructures/skiplist"
)

type MessageType rune

const (
	InsertMessage MessageType = 'I'
	QueryMessage  MessageType = 'Q'
)

const MessageLength = 9

type Message [MessageLength]byte

func MessageFromSlice(b []byte) Message {
	return Message(b)
}

func (msg *Message) String() string {
	return fmt.Sprintf("%c %d %d", rune(msg.Type()), msg.MinTime(), msg.MaxTime())
}

func (msg *Message) Type() MessageType {
	return MessageType(msg[0])
}

func (msg *Message) readInt32(offset uint) int32 {
	val := binary.BigEndian.Uint32(msg[offset:])
	return int32(val)
}

func (msg *Message) Timestamp() int32 {
	return msg.readInt32(1)
}

func (msg *Message) Price() int32 {
	return msg.readInt32(5)
}

func (msg *Message) MinTime() int32 {
	return msg.readInt32(1)
}

func (msg *Message) MaxTime() int32 {
	return msg.readInt32(5)
}

func computeMean(sl *skiplist.SkipList, start, end int) int {
	mean, n := 0, 0
	for _, value := range sl.RangeByScore(start, end) {
		mean += int(value.(int32))
		n++
	}

	if n == 0 {
		return 0
	}

	return mean / n
}

func connHandler(conn net.Conn) {
	defer log.Printf("%s disconnected\n", conn.RemoteAddr())
	defer conn.Close()

	sl := skiplist.NewSkipList(16)

	for {
		b := make([]byte, MessageLength)
		if _, err := io.ReadFull(conn, b); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			log.Fatalf("Read error: %s\n", err)
		}

		msg := MessageFromSlice(b)

		switch t := msg.Type(); t {
		case InsertMessage:
			start := int(msg.Timestamp())
			end := start
			existing := sl.RangeByScore(start, end)
			if len(existing) == 0 {
				sl.Insert(int(msg.Timestamp()), msg.Price())
			}
		case QueryMessage:
			mean := computeMean(sl, int(msg.MinTime()), int(msg.MaxTime()))
			outBuffer := make([]byte, 4)
			binary.BigEndian.PutUint32(outBuffer, uint32(mean))
			if _, err := conn.Write(outBuffer); err != nil {
				log.Fatalf("Write error: %s\n", err)
			}
		default:
			// Invalid message
			return
		}
	}
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
