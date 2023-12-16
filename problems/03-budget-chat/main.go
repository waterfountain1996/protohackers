package main

import (
	"io"
	"log"

	"github.com/waterfountain1996/protohackers/problems/03-budget-chat/client"
	"github.com/waterfountain1996/protohackers/problems/03-budget-chat/server"
)

const MaxUsernameLength = 16

const UsernamePrompt = "Welcome to budgetchat! What shall I call you?\n"

const InvalidUsernameMessage = "Usernames must only contain ASCII letters, digits or underscores\n"

func isUsernameValid(username []byte) bool {
	if len(username) == 0 || len(username) > MaxUsernameLength {
		return false
	}

	for _, b := range username {
		if !((b >= 'A' && b <= 'z') || (b >= '0' && b <= '9') || b == '_') {
			return false
		}
	}

	return true
}

func handleClient(c *client.Client, messageChan chan<- *client.Message, notify func()) error {
	if err := c.Write([]byte(UsernamePrompt)); err != nil {
		return err
	}

	username, _ := c.Read()
	if !isUsernameValid(username) {
		log.Printf("Client %v provided an invald username\n", c)
		if err := c.Write([]byte(InvalidUsernameMessage)); err != nil {
			return err
		}

		return nil
	}

	c.Username = string(username)

	// Notify user joined
	notify()

	for {
		text, err := c.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("Read error: %s\n", err)
		}

		msg := client.NewMessage(string(text), c)
		messageChan <- msg
	}

	return nil
}

func main() {
	server := server.NewServer()

	if err := server.RunForever(":10000", handleClient); err != nil {
		log.Fatal(err)
	}
}
