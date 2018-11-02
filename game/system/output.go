package system

import (
	"aragno/ecs"
	"aragno/game/component"
	"reflect"
	"unsafe"
)

// StateOutputSystem serializes game state and outputs to listeners
type StateOutputSystem struct {
	aether    *ecs.Aether
	hermes    *ecs.Hermes
	registrar *ecs.EntityRegistrar
}

// Register registers ecs internals
func (sos *StateOutputSystem) Register(aether *ecs.Aether, hermes *ecs.Hermes, registrar *ecs.EntityRegistrar) {
	sos.aether = aether
	sos.hermes = hermes
	sos.registrar = registrar
}

// NewStateOutputSystem constructor
func NewStateOutputSystem() *StateOutputSystem {
	return &StateOutputSystem{}
}

// Update serializes game state and sends it off to websocket writers, also sets "owner" for each websocket
func (sos *StateOutputSystem) Update(dt float64) {
	// Build game state
	state := component.GameState{Tid: component.GameStateTid}
	sos.SerializeBodies(&state)
	sos.SerializeLegs(&state)

	// Send out game state to players
	for _, v := range sos.aether.RetrieveType(reflect.TypeOf(&component.PlayerInput{})) {
		vcast := v.(*component.PlayerInput)
		if vcast.Disconnected {
			continue
		}
		state.OwnerId = uintptr(unsafe.Pointer(vcast.Conn))
		vcast.OutputChan <- state
	}
}

// SerializeBodies serializes spider body game states
func (sos *StateOutputSystem) SerializeBodies(state *component.GameState) {
	for id, v := range sos.aether.RetrieveType(reflect.TypeOf(&component.SpiderBody{})) {
		bs := component.BodyState{}
		bs.Name = v.(*component.SpiderBody).Name
		if v, err := sos.aether.Retrieve(id, reflect.TypeOf(&component.Pose{})); err == nil {
			bs.Pose = *(v.(*component.Pose))
		}
		state.Bodies = append(state.Bodies, bs)
	}
}

// SerializeLegs serializes spider leg game states
func (sos *StateOutputSystem) SerializeLegs(state *component.GameState) {
	for id, _ := range sos.aether.RetrieveType(reflect.TypeOf(&component.SpiderLeg{})) {
		ls := component.LegState{}

		if o, err := sos.aether.Retrieve(id, reflect.TypeOf(&component.Owner{})); err == nil {
			ownerId := o.(*component.Owner).Id
			if i, err := sos.aether.Retrieve(ownerId, reflect.TypeOf(&component.PlayerInput{})); err == nil {
				ls.ConnId = uintptr(unsafe.Pointer((i.(*component.PlayerInput).Conn)))
			}
		}

		if p, err := sos.aether.Retrieve(id, reflect.TypeOf(&component.Pose{})); err == nil {
			ls.Pose = *(p.(*component.Pose))
		}
		state.Legs = append(state.Legs, ls)
	}
}
