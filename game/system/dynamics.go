package system

import (
	"aragno/ecs"
	"aragno/game/component"
	"aragno/dynamo"
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
	// get all the bodies

	for id, bodyUncast := range ms.aether.RetrieveType(component.BodyType) {
//		fmt.Println(id)
		body := bodyUncast.(*dynamo.Body)
		// input velocites from the players if this is a player type
//		if playerInputUncast, err := ms.aether.Retrieve(id, component.PlayerInputType); err == nil {
			//pi := playerInputUncast.(*component.PlayerInput)

//		}
		// if the thing is already moving
		if velocityUncast, err := ms.aether.Retrieve(id, component.VelocityType); err == nil {
			if posUncast, err := ms.aether.Retrieve(id, component.PoseType); err == nil {
				pos := posUncast.(*component.Pose) 
				vel := velocityUncast.(*component.Velocity) 
				dynamo.UpdateBody(body,zero.Pose{0,0,0},dt)	
				body.Shapes[0].SetPosition(func(x float64, y float64, z float64 ){
					pos.X = x
					pos.Y= y
				})
				body.Shapes[0].SetVelocity(func(x float64, y float64, z float64 ){
					vel.Dx = x
					vel.Dy= y
				})

			}
		}
	}
}
