package game

import (
	"aragno/ecs"
	"github.com/gorilla/websocket"
)

// Pose pose for a rigid body component
type Pose struct {
	X     float32
	Y     float32
	Theta float32
}

// PivotRoot indicates a connection to another rigid body with a pivot joint
type PivotRoot struct {
	Root ecs.EntityId
	X    float32
	Y    float32
}

// Serializable tag component for indicating an entity should be serialized
type Serializable struct {
}

// Gametype component to name what kind of object something is for serialization
type GameObjectType struct {
	Type string
}

// Player the owning connection and input for this entity
type PlayerInput struct {
	X            float32
	Y            float32
	Valid        bool
	Clicked      bool
	Disconnected bool
	OutputChan   chan GameState
	Conn         *websocket.Conn
}

// StateOutput the state of the game sent to the clients
type GameState struct {
	message string
}
