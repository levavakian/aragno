package component

import (
	"github.com/gorilla/websocket"
	"reflect"
)

// Player the owning connection and input for this entity
type PlayerInput struct {
	X            float64
	Y            float64
	Valid        bool
	Clicked      bool
	Disconnected bool
	OutputChan   chan GameState
	Conn         *websocket.Conn
}

var PlayerInputType = reflect.TypeOf(&PlayerInput{})
