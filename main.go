package main

import (
	"aragno/game"
	"aragno/io"
	"net/http"
	"fmt"
)

func main() {
	input := make(chan game.PlayerInput)
	go game.Loop(input)
	http.HandleFunc("/connect", io.Connect(input))
	fmt.Println(http.ListenAndServe("localhost:8000", nil))
}
