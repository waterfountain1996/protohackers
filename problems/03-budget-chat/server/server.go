package server

import (
	"log"
	"net"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/waterfountain1996/protohackers/problems/03-budget-chat/client"
)

type Server struct {
	clients    []*client.Client
	clientLock sync.RWMutex

	// Channel for relaying messages between clients
	messageChan chan *client.Message
}

func NewServer() *Server {
	return &Server{
		clients:     []*client.Client{},
		messageChan: make(chan *client.Message),
	}
}

// Add client to the list of currently connected ones
func (srv *Server) addClient(c *client.Client) {
	srv.clientLock.Lock()
	defer srv.clientLock.Unlock()

	srv.clients = append(srv.clients, c)
}

// Remove client from the list of currently connected ones
// and inform others that the user has left.
func (srv *Server) removeClient(c *client.Client) {
	if c.Joined() {
		defer srv.announceMembership(c, true /* leaving */)
	}

	srv.clientLock.Lock()
	defer srv.clientLock.Unlock()

	idx := slices.Index(srv.clients, c)
	srv.clients = slices.Delete(srv.clients, idx, idx+1)
}

func (srv *Server) activeUsers() []*client.Client {
	srv.clientLock.RLock()
	defer srv.clientLock.RUnlock()

	clients := []*client.Client{}
	for _, c := range srv.clients {
		if c.Joined() {
			clients = append(clients, c)
		}
	}

	return clients
}

// Write `data` to all clients except the one passed in `exception`.
func (srv *Server) writeToAllExcept(data []byte, exception *client.Client) {
	for _, recepient := range srv.activeUsers() {
		if recepient == exception {
			continue
		}

		go func(c *client.Client) {
			if err := c.Write(data); err != nil {
				if !os.IsTimeout(err) {
					log.Printf("ERROR: Write error: %s\n", err)
				}
			}
		}(recepient)
	}
}

func (srv *Server) announceMembership(c *client.Client, leaving bool) {
	verb := "entered"
	if leaving {
		verb = "left"
	}

	var b strings.Builder
	b.WriteString("* ")
	b.WriteString(c.Username)
	b.WriteString(" has ")
	b.WriteString(verb)
	b.WriteString(" the room")

	srv.writeToAllExcept(append([]byte(b.String()), '\n'), c)
}

// Send a list of active users to a newly connected client.
func (srv *Server) listMembers(c *client.Client) {
	names := []string{}
	for _, user := range srv.activeUsers() {
		if user == c {
			continue
		}
		names = append(names, user.Username)
	}

	var b strings.Builder
	b.WriteString("* This room ")

	if len(names) > 0 {
		b.WriteString("contains: ")
		b.WriteString(strings.Join(names, ", "))
	} else {
		b.WriteString("is empty")
	}

	data := append([]byte(b.String()), '\n')
	if err := c.Write(data); err != nil {
		if !os.IsTimeout(err) {
			log.Printf("ERROR: Write error: %s\n", err)
		}
	}
}

// Format message according to the protocol specification.
func messageToString(msg *client.Message) string {
	var b strings.Builder
	b.WriteRune('[')
	b.WriteString(msg.Sender.Username)
	b.WriteString("] ")
	b.WriteString(msg.Text)
	return b.String()
}

// Read messages from `srv.messageChan` and relay to all joined clients.
func (srv *Server) relayMessages() {
	for msg := range srv.messageChan {
		data := append([]byte(messageToString(msg)), '\n')
		srv.writeToAllExcept(data, msg.Sender)
	}
}

type ClientHandler func(*client.Client, chan<- *client.Message, func()) error

// Run the chat server on the given 'address'.
func (srv *Server) RunForever(address string, handler ClientHandler) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	go srv.relayMessages()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		log.Printf("New TCP connection: %s\n", conn.RemoteAddr())

		c := client.NewClient(conn)
		srv.addClient(c)

		go func(c *client.Client) {
			defer conn.Close()
			defer srv.removeClient(c)

			handler(c, srv.messageChan, func() {
				srv.listMembers(c)
				srv.announceMembership(c, false)
			})

			log.Printf("%s disconnected\n", conn.RemoteAddr())
		}(c)
	}
}
