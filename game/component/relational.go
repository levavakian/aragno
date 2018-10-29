package component

import (
	"aragno/ecs"
	"reflect"
)

// Owner indicates a player owns this object
type Owner struct {
	Id ecs.EntityId
}

var OwnerType = reflect.TypeOf(&Owner{})

// Children indicates that the entities pointed to by this parent share a lifetime with parent entity
type Children struct {
	Ids []ecs.EntityId
}

var ChildrenType = reflect.TypeOf(&Children{})
