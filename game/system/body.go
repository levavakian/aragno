package system

import (
	"aragno/ecs"
	"aragno/game/component"
)

// BodySystem controls the creation of spider bodies
type BodySystem struct {
	aether    *ecs.Aether
	hermes    *ecs.Hermes
	registrar *ecs.EntityRegistrar
}

// Register registers ecs internals
func (bs *BodySystem) Register(aether *ecs.Aether, hermes *ecs.Hermes, registrar *ecs.EntityRegistrar) {
	bs.aether = aether
	bs.hermes = hermes
	bs.registrar = registrar
}

// Init initializes spider bodies
func (bs *BodySystem) Init() {
	spiders := [...]string{"Carl", "Frog"}
	for _, name := range spiders {
		id := bs.registrar.NewId()
		bs.aether.Register(id, &component.SpiderBody{Name: name})
		bs.aether.Register(id, &component.Pose{5, 5, 5})
		bs.hermes.Send(&ecs.Message{Pipe: EntityCreatedPipe, EntityId: id})
	}
}

// NewBodySystem constructor
func NewBodySystem() *BodySystem {
	return &BodySystem{}
}

// Update does nothing for now
func (bs *BodySystem) Update(dt float64) {
	// TODO
}
