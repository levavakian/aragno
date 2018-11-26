package component

import (
	"aragno/dynamo"
	"reflect"
)

// Pose pose for a rigid body component
type Pose struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Theta float64 `json:"theta"`
}

var PoseType = reflect.TypeOf(&Pose{})

// Velocity velocity for movement frame to frame
type Velocity struct {
	Dx     float64
	Dy     float64
	Dtheta float64
}

var VelocityType = reflect.TypeOf(&Velocity{})

// The body struct is temporarily implemented in the dynamo library
var BodyType = reflect.TypeOf(&dynamo.Body{})
