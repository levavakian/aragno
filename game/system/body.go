package system

import (
	"aragno/ecs"
	"aragno/game/component"
)

type BodySystem struct {
	aether    *ecs.Aether
	hermes    *ecs.Hermes
	registrar *ecs.EntityRegistrar
}

func (bs *BodySystem) Register(aether *ecs.Aether, hermes *ecs.Hermes, registrar *ecs.EntityRegistrar) {
	bs.aether = aether
	bs.hermes = hermes
	bs.registrar = registrar
}

func (bs *BodySystem) Init() {
	spiders := [...]string{"Carl", "Frog"}
	for _, name := range spiders {
		id := bs.registrar.NewId()
		bs.aether.Register(id, &component.SpiderBody{Name: name})
		bs.aether.Register(id, &component.Pose{5, 5, 5})
		bs.hermes.Send(&ecs.Message{Pipe: EntityCreatedPipe, EntityId: id})
	}
}

func NewBodySystem() *BodySystem {
	return &BodySystem{}
}

func (bs *BodySystem) Update(dt float32) {
	// TODO
}
