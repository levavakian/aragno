package system

import (
	"aragno/dynamo"
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
	// Starting positions for spiders
	cid := ms.registrar.NewId()
	ms.aether.Register(cid, &component.SpiderBody{Name: "Carl"})
	ms.aether.Register(cid, &component.Pose{300, 50, -0.78539816339})
	ms.aether.Register(cid, &component.Velocity{0, 0, 0})
	ms.aether.Register(cid, dynamo.NewRectangleBody("gimme a name", 300, 50, 1, zero.Pose{300, 50, 0.78}))
	ms.hermes.Send(&ecs.Message{Pipe: EntityCreatedPipe, EntityId: cid})

	fid := ms.registrar.NewId()
	ms.aether.Register(fid, &component.SpiderBody{Name: "Frog"})
	ms.aether.Register(fid, &component.Pose{-300, 50, 0.78539816339})
	ms.aether.Register(fid, &component.Velocity{0, 0, 0})
	ms.aether.Register(fid, dynamo.NewRectangleBody("gimme a name", 300, 50, 1, zero.Pose{-300, 50, 0.78}))
	ms.hermes.Send(&ecs.Message{Pipe: EntityCreatedPipe, EntityId: fid})

	// TODO: load from config
	ms.Surfaces = append(ms.Surfaces, zero.Rectangle{zero.Point{-500, -150}, zero.Point{500, -150}, zero.Point{500, -250}})

	// Register callback for sending out maps to subscribers
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
