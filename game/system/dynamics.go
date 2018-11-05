package system

import (
	"aragno/ecs"
	"aragno/game/component"
	"aragno/zero"
)

// Dynamics move players and decide velocities based on the forces on a player
type DynamicsSystem struct {
	aether    *ecs.Aether
	hermes    *ecs.Hermes
	registrar *ecs.EntityRegistrar
}

// Register registers ecs internals
func (ms *DynamicsSystem) Register(aether *ecs.Aether, hermes *ecs.Hermes, registrar *ecs.EntityRegistrar) {
	ms.aether = aether
	ms.hermes = hermes
	ms.registrar = registrar
}

// NewDynamicsSystem ecs internals
func NewDynamicsSystem() *DynamicsSystem {
	return &DynamicsSystem{}
}

// Update updates poses according to their velocities
func (ms *DynamicsSystem) Update(dt float64) {
	for id, velUncast := range ms.aether.RetrieveType(component.VelocityType) {
		vel := velUncast.(*component.Velocity)
		if posUncast, err := ms.aether.Retrieve(id, component.PoseType); err == nil {
			pos := posUncast.(*component.Pose)
			pos.X = pos.X + vel.Dx*dt
			pos.Y = pos.Y + vel.Dy*dt
			pos.Theta = zero.NormalizeBipolar(pos.Theta + vel.Dtheta*dt)
		}
	}
}
