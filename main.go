package main

import (
	"aragno/game"
	"aragno/game/component"
	"aragno/io"
	"fmt"
	"net/http"
)

func main() {
	// Communication channels
	input := make(chan component.PlayerInput)

	// Game loop
	fmt.Println("Initiating game loop")
	go game.Loop(input)

	// Server
	http.HandleFunc("/connect", io.Connect(input))
	fmt.Println(http.ListenAndServe("localhost:8000", nil))
}
