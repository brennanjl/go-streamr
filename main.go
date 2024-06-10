package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

const streamID = "streams.dimo.eth/firehose/weather"
const apiKey = "OWZjODdlN2VjNmNiNGMzYTgzNjRmZmExNzYwNmUxN2Y"

func main() {

	path := "/streams/" + url.QueryEscape(streamID) + "/subscribe"
	// Define the WebSocket server URL
	u := url.URL{Scheme: "ws", Host: "localhost:7170", Path: path, RawQuery: "apiKey=" + apiKey}
	log.Printf("connecting to %s", u.String())

	// Create a WebSocket connection
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	// Listen for incoming messages
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				return
			}
			fmt.Printf("recv: %s", message)
		}
	}()

	// Setup a channel to handle interrupts
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
}
