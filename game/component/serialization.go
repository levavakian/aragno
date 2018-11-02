package component

import (
	"aragno/zero"
	"reflect"
)

const (
	GameStateTid = "GameState"
	MapStateTid  = "MapState"
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
	Tid     string      `json:"tid"`
	Bodies  []BodyState `json:"bodies"`
	Legs    []LegState  `json:"legs"`
	OwnerId uintptr     `json:"owner_id"`
}

// MapState the state of the map
type MapState struct {
	Tid      string           `json:"tid"`
	Surfaces []zero.Rectangle `json:"surfaces"`
}
