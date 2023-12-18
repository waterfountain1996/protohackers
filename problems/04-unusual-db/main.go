package main

import (
	"log"
	"net"
	"strings"
	"sync"
)

const Version string = "KeeValue Store 6.9"

type DB struct {
	items map[string]string
	lock  sync.RWMutex
}

func NewDB() *DB {
	return &DB{
		items: make(map[string]string),
	}
}

func (db *DB) Get(key string) string {
	db.lock.RLock()
	defer db.lock.RUnlock()

	return db.items[key]
}

func (db *DB) Set(key string, value string) {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.items[key] = value
}

func handleRequest(db *DB, pc net.PacketConn, addr net.Addr, data []byte) {
	// Max request length
	if len(data) > 1000 {
		return
	}

	key, value, isInsert := strings.Cut(string(data), "=")
	if isInsert {
		if key != "version" {
			db.Set(key, value)
		}
	} else {
		var result string

		if key == "version" {
			result = Version
		} else {
			result = db.Get(key)
		}

		payload := strings.Join([]string{key, result}, "=")
		if _, err := pc.WriteTo([]byte(payload), addr); err != nil {
			log.Printf("UDP write error: %s\n", err)
		}
	}
}

func main() {
	pc, err := net.ListenPacket("udp", ":10000")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	db := NewDB()

	for {
		b := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(b)
		if err != nil {
			log.Printf("UDP read error: %s\n", err)
			continue
		}

		log.Printf("Received %d byte long packet from %s\n", n, addr.String())

		go handleRequest(db, pc, addr, b[:n])
	}
}
