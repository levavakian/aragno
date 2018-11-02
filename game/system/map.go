package system

import (
	"aragno/ecs"
	"aragno/game/component"
	"aragno/zero"
)

type MapSystem struct {
	aether    *ecs.Aether
	hermes    *ecs.Hermes
	registrar *ecs.EntityRegistrar
	Surfaces  []zero.Rectangle
}

// Register registers ecs internals
func (ms *MapSystem) Register(aether *ecs.Aether, hermes *ecs.Hermes, registrar *ecs.EntityRegistrar) {
	ms.aether = aether
	ms.hermes = hermes
	ms.registrar = registrar
}

// Init initializes spider bodies
func (ms *MapSystem) Init() {
	// TODO: load from config
	ms.Surfaces = append(ms.Surfaces, zero.Rectangle{zero.Point{-10, 1}, zero.Point{10, 1}, zero.Point{10, -1}})

	ms.hermes.AddCallback(PlayerCreatedPipe, func(msg *ecs.Message) {
		if input, err := ms.aether.Retrieve(msg.EntityId, component.PlayerInputType); err == nil {
			// Do deep copy of surfaces and transmit
			surfaces := []zero.Rectangle{}
			for _, rec := range ms.Surfaces {
				surfaces = append(surfaces, rec)
			}
			mapState := component.MapState{Tid: component.MapStateTid, Surfaces: surfaces}

			input.(*component.PlayerInput).OutputChan <- mapState
		}
	})
}

func NewMapSystem() *MapSystem {
	return &MapSystem{}
}

// Update does nothing for now
func (ms *MapSystem) Update(dt float64) {
	// TODO
}
