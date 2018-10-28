package main

import (
	"aragno/game"
	"aragno/io"
	"fmt"
	"net/http"
)

func main() {
	// Communication channels
	input := make(chan game.PlayerInput)

	// Game loop
	go game.Loop(input)

	// Server
	http.HandleFunc("/connect", io.Connect(input))
	fmt.Println(http.ListenAndServe("localhost:8000", nil))
}
