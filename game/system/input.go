package system

import (
	"aragno/ecs"
	"aragno/game/component"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

// PlayerInputSystem handles updating player inputs and creating new players
type PlayerInputSystem struct {
	aether    *ecs.Aether
	hermes    *ecs.Hermes
	registrar *ecs.EntityRegistrar
	InputChan chan component.PlayerInput
	Inputs    map[*websocket.Conn]component.PlayerInput
	CloseSig  chan struct{}
	mux       sync.Mutex
}

// NewPlayerInputSystem constructor
func NewPlayerInputSystem(inputChan chan component.PlayerInput) *PlayerInputSystem {
	return &PlayerInputSystem{InputChan: inputChan, Inputs: make(map[*websocket.Conn]component.PlayerInput), CloseSig: make(chan struct{})}
}

// Register registers ecs internals
func (pis *PlayerInputSystem) Register(aether *ecs.Aether, hermes *ecs.Hermes, registrar *ecs.EntityRegistrar) {
	pis.aether = aether
	pis.hermes = hermes
	pis.registrar = registrar
}

// Init launches input collector
func (pis *PlayerInputSystem) Init() {
	go pis.Collect()
}

// Update updates inputs, creates new players if new input is detected, destroys player if disconnected
func (pis *PlayerInputSystem) Update(dt float64) {
	inputs := pis.PopInputs()

	// Create mapping of Conn -> Entity
	connToEntity := make(map[*websocket.Conn]ecs.EntityId)
	for k, v := range pis.aether.RetrieveType(component.PlayerInputType) {
		connToEntity[v.(*component.PlayerInput).Conn] = k
	}

	// Create
	for conn, input := range inputs {
		playerId, exists := connToEntity[conn]
		if !exists {
			// Create player
			playerId = pis.registrar.NewId()

			// Create leg
			legId := pis.CreateLeg(playerId)

			// Register player components
			pis.aether.Register(playerId, &input)
			pis.aether.Register(playerId, &component.Children{[]ecs.EntityId{legId}})

			// Emit player and leg creation events
			pis.hermes.Send(&ecs.Message{Pipe: EntityCreatedPipe, EntityId: playerId})
			pis.hermes.Send(&ecs.Message{Pipe: PlayerCreatedPipe, EntityId: playerId})
			pis.hermes.Send(&ecs.Message{Pipe: EntityCreatedPipe, EntityId: legId})

		} else {
			pis.aether.Register(playerId, &input)
		}

		if input.Disconnected {
			close(input.OutputChan)
			pis.hermes.Send(&ecs.Message{Pipe: EntityDestroyedPipe, EntityId: playerId})
		}
	}
}

// CreateLeg creates a new spider leg
func (pis *PlayerInputSystem) CreateLeg(playerId ecs.EntityId) ecs.EntityId {
	id := pis.registrar.NewId()
	pis.aether.Register(id, &component.Owner{playerId})
	pis.aether.Register(id, &component.SpiderLeg{})
	pis.aether.Register(id, &component.Velocity{5, 5, 5})
	pis.aether.Register(id, &component.Pose{5, 5, 5})

	return id
}

// Cleanup shuts down collector
func (pis *PlayerInputSystem) Cleanup() {
	close(pis.CloseSig)
}

// Collect collects inputs received from ws readers and stores them for retrieval
func (pis *PlayerInputSystem) Collect() {
	for {
		select {
		case input := <-pis.InputChan:
			pis.ReceiveInput(input)
		case _, ok := <-pis.CloseSig:
			if !ok {
				fmt.Println("Shutting down input collecter")
				break
			}
		}
	}
}

// ReceiveInput does a thread safe update on the map of received player inputs
func (pis *PlayerInputSystem) ReceiveInput(input component.PlayerInput) {
	pis.mux.Lock()
	defer pis.mux.Unlock()
	pis.Inputs[input.Conn] = input
}

// PopInputs does a thread safe retrieve of player inputs since last Pop
func (pis *PlayerInputSystem) PopInputs() map[*websocket.Conn]component.PlayerInput {
	pis.mux.Lock()
	defer pis.mux.Unlock()

	inputs := pis.Inputs
	pis.Inputs = make(map[*websocket.Conn]component.PlayerInput)
	return inputs
}
