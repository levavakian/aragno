package component

import (
	"aragno/ecs"
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

// PivotRoot indicates a connection to another rigid body with a pivot joint
type PivotRoot struct {
	Root ecs.EntityId
	X    float64
	Y    float64
}

var PivotRootType = reflect.TypeOf(&PivotRoot{})
