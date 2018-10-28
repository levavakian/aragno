package game

import (
	"aragno/ecs"
	"reflect"
)

type StateOutputSystem struct {
	aether    *ecs.Aether
	hermes    *ecs.Hermes
	registrar *ecs.EntityRegistrar
}

func (sos *StateOutputSystem) Register(aether *ecs.Aether, hermes *ecs.Hermes, registrar *ecs.EntityRegistrar) {
	sos.aether = aether
	sos.hermes = hermes
	sos.registrar = registrar
}

func NewStateOutputSystem() *StateOutputSystem {
	return &StateOutputSystem{}
}

func (sos *StateOutputSystem) Update(dt float32) {
	stored, err := sos.aether.RetrieveType(reflect.TypeOf(&PlayerInput{}))
	if err != nil {
		panic(err)
	}
	for _, v := range stored {
		vcast := v.(*PlayerInput)
		if vcast.Disconnected {
			continue
		}
		vcast.OutputChan <- GameState{message: "hello"}
	}
}
