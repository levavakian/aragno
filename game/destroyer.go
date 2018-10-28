package game

import (
	"aragno/ecs"
	"fmt"
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
		fmt.Printf("EntityDestroyerSystem.Update: destroyed entity %s\n", entity)
		eds.aether.DeregisterAll(entity)
	}
	eds.destroyed = []ecs.EntityId{}
}
