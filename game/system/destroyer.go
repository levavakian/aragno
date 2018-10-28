package system

import (
	"aragno/ecs"
	"aragno/game/component"
	"reflect"
)

type EntityDestroyerSystem struct {
	aether    *ecs.Aether
	hermes    *ecs.Hermes
	registrar *ecs.EntityRegistrar
	destroyed []ecs.EntityId
}

func (eds *EntityDestroyerSystem) Register(aether *ecs.Aether, hermes *ecs.Hermes, registrar *ecs.EntityRegistrar) {
	eds.aether = aether
	eds.hermes = hermes
	eds.registrar = registrar
}

func (eds *EntityDestroyerSystem) Init() {
	eds.hermes.AddCallback(EntityDestroyedPipe, func(msg *ecs.Message) {
		eds.destroyed = append(eds.destroyed, msg.EntityId)
	})
}

func NewEntityDestroyerSystem() *EntityDestroyerSystem {
	return &EntityDestroyerSystem{}
}

func (eds *EntityDestroyerSystem) Update(dt float32) {
	for _, entity := range eds.destroyed {
		eds.Destroy(entity)
	}
	eds.destroyed = []ecs.EntityId{}
}

func (eds *EntityDestroyerSystem) Destroy(id ecs.EntityId) {
	// Destroy children frist
	if children, err := eds.aether.Retrieve(id, reflect.TypeOf(&component.Children{})); err == nil {
		for _, childId := range children.(*component.Children).Ids {
			eds.Destroy(childId)
		}
	}
	eds.aether.DeregisterAll(id)
}
