package system

import (
	"aragno/ecs"
	"aragno/game/component"
	"aragno/zero"
)

// MoveSystem updates poses based on their velocities
type MoveSystem struct {
	aether    *ecs.Aether
	hermes    *ecs.Hermes
	registrar *ecs.EntityRegistrar
}

// Register registers ecs internals
func (ms *MoveSystem) Register(aether *ecs.Aether, hermes *ecs.Hermes, registrar *ecs.EntityRegistrar) {
	ms.aether = aether
	ms.hermes = hermes
	ms.registrar = registrar
}

// NewMoveSystem ecs internals
func NewMoveSystem() *MoveSystem {
	return &MoveSystem{}
}

// Update updates poses according to their velocities
func (ms *MoveSystem) Update(dt float64) {
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
