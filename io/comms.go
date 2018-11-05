package io

import (
	"aragno/game/component"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

// Read starts reading from a websocket connection and publishing to input channel
func Read(connection *websocket.Conn, inputChan chan component.PlayerInput, outputChan chan interface{}) {
	fmt.Printf("Starting reader for %s\n", &connection)

	// Send initial message to indicate connection
	inputChan <- component.PlayerInput{X: 0, Y: 0, Valid: false, Clicked: false, Disconnected: false, Conn: connection, OutputChan: outputChan}

	// Read until connection is closed
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			fmt.Printf("Read connection closed for %s\n", &connection)
			connection.Close()
			inputChan <- component.PlayerInput{Conn: connection, Disconnected: true, OutputChan: outputChan}
			break
		}
		fmt.Printf("%s\n", message)
	}
}

// Write starts writing to websocket connection when game state is received from the game loop
func Write(connection *websocket.Conn, outputChan chan interface{}) {
	fmt.Printf("Starting writer for %s\n", &connection)
	for {
		// Receive game state information
		output, ok := (<-outputChan)

		// Close out if game loop has terminated the write connection
		if !ok {
			fmt.Printf("Close requested for %s\n", &connection)
			break
		}

		// Publish game state
		err := connection.WriteJSON(output)
		if err != nil {
			fmt.Printf("Write connection failed with %s for %s\n", err, &connection)
			connection.Close()
		}
	}
}

// Connect upgrades websocket connection and starts player input ws reader and game state ws writer
func Connect(inputChan chan component.PlayerInput) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, req *http.Request) {
		// TODO: handle COR correctly
		upgrader := websocket.Upgrader{}
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}

		// Upgrade http connection and generate new output channel for writing
		connection, err := upgrader.Upgrade(writer, req, nil)
		if err != nil {
			fmt.Printf("Connection failed: %s\n", err)
			return
		}
		outputChan := make(chan interface{})

		// Start reader and writer loops
		go Read(connection, inputChan, outputChan)
		go Write(connection, outputChan)
	}
}
