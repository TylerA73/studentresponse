package main

import (
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

func startSocket() *gosocketio.Server {
	srv := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	type Message struct {
		Room   string `json:"room"`
		String string `json:"message"`
	}

	srv.On("message", func(c *gosocketio.Channel, msg Message) {
		// Redirect messages to their class "rooms".
		c.BroadcastTo(msg.Room, "message", msg.String)
	})

	return srv
}
