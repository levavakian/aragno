package component

import (
	"reflect"
)

// Serializable tag component for indicating an entity should be serialized
type Serializable struct {
}

var SerializableType = reflect.TypeOf(&Serializable{})

// BodyState the state of a spider body
type BodyState struct {
	Name string `json:"name"`
	Pose Pose   `json:"pose"`
}

// LegState the state of a spider leg
type LegState struct {
	Pose   Pose    `json:"pose"`
	ConnId uintptr `json:"conn_id"`
}

// GameState overall game state
type GameState struct {
	Bodies  []BodyState `json:"bodies"`
	Legs    []LegState  `json:"legs"`
	OwnerId uintptr     `json:"owner_id"`
}
