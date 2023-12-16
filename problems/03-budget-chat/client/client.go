package client

import (
	"bufio"
	"io"
	"net"
	"time"
)

const (
	clientWriteTimeout = time.Second * 5
)

type Client struct {
	// Client's underlying TCP connection.
	conn net.Conn

	// Protocol specifies that all messages are '\n' terminated
	// so we can use a scanner for reading.
	scanner *bufio.Scanner

	// Client's username. Empty string indicates that they have not joined yet.
	Username string
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:    conn,
		scanner: bufio.NewScanner(conn),
	}
}

func (c *Client) Joined() bool {
	return c.Username != ""
}

func (c *Client) Write(data []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(clientWriteTimeout))
	_, err := c.conn.Write(data)
	return err
}

func (c *Client) Read() ([]byte, error) {
	if !c.scanner.Scan() {
		err := c.scanner.Err()
		if err == nil {
			err = io.EOF
		}

		return nil, err
	}
	return c.scanner.Bytes(), nil
}

type Message struct {
	Text   string
	Sender *Client
}

func NewMessage(text string, sender *Client) *Message {
	return &Message{
		Text:   text,
		Sender: sender,
	}
}
