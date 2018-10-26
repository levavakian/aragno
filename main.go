package main

import (
	"aragno/io"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func gameLoop(collector *io.CollectedInputs) {
	fmt.Println("Starting game loop")

	// Set up output dispatcher
	outputDispatch := io.OutputDispatch{OutputChans: make(map[*websocket.Conn]chan io.StateOutput)}

	// Game loop
	for {
		// Collect and consume inputs
		collectedInputs := collector.Pop()
		outputDispatch.ConsumeInputs(collectedInputs)

		// Do game stuff
		state := io.StateOutput{Body: io.BodyState{X: 1, Y: 1, Theta: 1}, Legs: []io.LegState{io.LegState{X: 2, Y: 2, Theta: 2}, io.LegState{X: 3, Y: 3, Theta: 3}}}

		// Output game state to websockets
		outputDispatch.Dispatch(state)

		// Server run rate at ~30Hz
		time.Sleep(time.Millisecond * 33)
	}
}

func main() {
	// Set up collecter to grab outputs
	collector := &io.CollectedInputs{InputChan: make(chan io.PlayerInput), Inputs: make(map[*websocket.Conn]io.PlayerInput), CloseSig: make(chan struct{})}
	go collector.Collect()

	go gameLoop(collector)

	http.HandleFunc("/connect", io.Connect(collector.InputChan))
	fmt.Println(http.ListenAndServe("localhost:8000", nil))
}
