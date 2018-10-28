package component

import (
	"aragno/ecs"
	"github.com/gorilla/websocket"
)

// Pose pose for a rigid body component
type Pose struct {
	X     float32 `json:"x"`
	Y     float32 `json:"y"`
	Theta float32 `json:"theta"`
}

// Owner indicates a player owns this object
type Owner struct {
	Id ecs.EntityId
}

// Children indicates that the entities pointed to by this parent share a lifetime with parent entity
type Children struct {
	Ids []ecs.EntityId
}

// PivotRoot indicates a connection to another rigid body with a pivot joint
type PivotRoot struct {
	Root ecs.EntityId
	X    float32
	Y    float32
}

// SpiderBody tags an entity as a spider body
type SpiderBody struct {
	Name string
}

// SpiderLeg tags an entity as a spider leg
type SpiderLeg struct {
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
type BodyState struct {
	Name string `json:"name"`
	Pose Pose   `json:"pose"`
}

type LegState struct {
	Pose   Pose    `json:"pose"`
	ConnId uintptr `json:"conn_id"`
}

type GameState struct {
	Bodies  []BodyState `json:"bodies"`
	Legs    []LegState  `json:"legs"`
	OwnerId uintptr     `json:"owner_id"`
}
