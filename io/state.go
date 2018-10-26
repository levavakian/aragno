package io

import (
	"github.com/gorilla/websocket"
)

// BodyState state of the spider body
type BodyState struct {
	X     float32 `json:"x"`
	Y     float32 `json:"y"`
	Theta float32 `json:"theta"`
}

// LegState state of a spider leg
type LegState struct {
	X     float32         `json:"x"`
	Y     float32         `json:"y"`
	Theta float32         `json:"theta"`
	Owner bool            `json:"owner"`
	Conn  *websocket.Conn `json:"-"`
}

// StateOutput combined game state
type StateOutput struct {
	Body      BodyState  `json:"body"`
	Legs      []LegState `json:"legs"`
	PlayerIdx int        `json:"player_idx"`
}

// PlayerInput incoming client side state
type PlayerInput struct {
	X            float32
	Y            float32
	Valid        bool
	Clicked      bool
	Disconnected bool
	Conn         *websocket.Conn
	OutputChan   chan StateOutput
}
