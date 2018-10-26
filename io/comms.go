package io

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

// Read starts reading from a websocket connection and publishing to input channel
func Read(connection *websocket.Conn, inputChan chan PlayerInput, outputChan chan StateOutput) {
	fmt.Printf("Starting reader for %s", &connection)

	// Send initial message to indicate connection
	inputChan <- PlayerInput{X: 0, Y: 0, Valid: false, Clicked: false, Disconnected: false, Conn: connection, OutputChan: outputChan}

	// Read until connection is closed
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			fmt.Printf("Read connection closed for %s", &connection)
			connection.Close()
			inputChan <- PlayerInput{Conn: connection, Disconnected: true}
			break
		}
		fmt.Printf("%s\n", message)
	}
}

// Write starts writing to websocket connection when game state is received from the game loop
func Write(connection *websocket.Conn, outputChan chan StateOutput) {
	fmt.Printf("Starting writer for %s", &connection)
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

// Websocket upgrader
var upgrader = websocket.Upgrader{}

// Connect upgrades websocket connection and starts player input ws reader and game state ws writer
func Connect(inputChan chan PlayerInput) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, req *http.Request) {
		// Upgrade http connection and generate new output channel for writing
		connection, err := upgrader.Upgrade(writer, req, nil)
		if err != nil {
			fmt.Printf("Connection failed: %s\n", err)
			return
		}
		outputChan := make(chan StateOutput)

		// Start reader and writer loops
		go Read(connection, inputChan, outputChan)
		go Write(connection, outputChan)
	}
}

// CollectedInputs a sink for websocket connections to dump latest client side state information into
type CollectedInputs struct {
	InputChan chan PlayerInput
	Inputs    map[*websocket.Conn]PlayerInput
	CloseSig  chan struct{}
	mux       sync.Mutex
}

// Add does a thread safe update on the map of received player inputs
func (collector *CollectedInputs) Add(input PlayerInput) {
	collector.mux.Lock()
	defer collector.mux.Unlock()
	collector.Inputs[input.Conn] = input
}

// Pop does a thread safe retrieve of player inputs since last PopS
func (collector *CollectedInputs) Pop() map[*websocket.Conn]PlayerInput {
	collector.mux.Lock()
	defer collector.mux.Unlock()

	inputs := collector.Inputs
	collector.Inputs = make(map[*websocket.Conn]PlayerInput)
	return inputs
}

// Collect collects inputs received from ws readers and stores them for retrieval
func (collecter *CollectedInputs) Collect() {
	for {
		select {
		case input := <-collecter.InputChan:
			collecter.Add(input)
		case _, ok := <-collecter.CloseSig:
			if !ok {
				fmt.Println("Shutting down input collecter")
				break
			}
		}
	}
}

// OutputDispatch handles the creation and dispatch of game state to ws writers
type OutputDispatch struct {
	OutputChans map[*websocket.Conn]chan StateOutput
}

// ConsumeInputs tracks new ws writers to publish state to and cleans up closed connections when readers shut down
func (dispatch *OutputDispatch) ConsumeInputs(inputs map[*websocket.Conn]PlayerInput) map[*websocket.Conn]bool {
	connUpdates := make(map[*websocket.Conn]bool)
	for inConn, inInput := range inputs {
		fmt.Printf("Consuming input from %s\n", &inConn)

		_, exists := dispatch.OutputChans[inConn]
		if !exists {
			dispatch.OutputChans[inConn] = inInput.OutputChan
			connUpdates[inConn] = true
		}

		if inInput.Disconnected {
			close(dispatch.OutputChans[inConn])
			delete(dispatch.OutputChans, inConn)
			connUpdates[inConn] = false
		}
	}
	return connUpdates
}

// Dispatch sends game state to all tracked ws writers
func (dispatch *OutputDispatch) Dispatch(state StateOutput) {
	for _, outChan := range dispatch.OutputChans {
		outChan <- state
	}
}
