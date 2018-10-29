package system

import (
	"aragno/ecs"
	"aragno/game/component"
)

// Handles destruction of entities
type EntityDestroyerSystem struct {
	aether    *ecs.Aether
	hermes    *ecs.Hermes
	registrar *ecs.EntityRegistrar
	destroyed []ecs.EntityId
}

// Register registers ecs internals
func (eds *EntityDestroyerSystem) Register(aether *ecs.Aether, hermes *ecs.Hermes, registrar *ecs.EntityRegistrar) {
	eds.aether = aether
	eds.hermes = hermes
	eds.registrar = registrar
}

// Init registers callback for entity destruction
func (eds *EntityDestroyerSystem) Init() {
	eds.hermes.AddCallback(EntityDestroyedPipe, func(msg *ecs.Message) {
		eds.destroyed = append(eds.destroyed, msg.EntityId)
	})
}

// NewEntityDestroyerSystem constructor
func NewEntityDestroyerSystem() *EntityDestroyerSystem {
	return &EntityDestroyerSystem{}
}

// Update deregisters all components tied to entities whose destruction has been requested
func (eds *EntityDestroyerSystem) Update(dt float64) {
	for _, entity := range eds.destroyed {
		eds.Destroy(entity)
	}
	eds.destroyed = []ecs.EntityId{}
}

// Destroy recursively destroys entity and any related children
func (eds *EntityDestroyerSystem) Destroy(id ecs.EntityId) {
	// Destroy children frist
	if children, err := eds.aether.Retrieve(id, component.ChildrenType); err == nil {
		for _, childId := range children.(*component.Children).Ids {
			eds.Destroy(childId)
		}
	}
	eds.aether.DeregisterAll(id)
}
