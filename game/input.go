package game

import (
	"sync"
	"github.com/gorilla/websocket"
	"aragno/ecs"
	"reflect"
	"fmt"
)

type PlayerInputSystem struct {
	aether    *ecs.Aether
	hermes    *ecs.Hermes
	registrar *ecs.EntityRegistrar
	InputChan chan PlayerInput
	Inputs    map[*websocket.Conn]PlayerInput
	CloseSig  chan struct{}
	mux       sync.Mutex
}

func NewPlayerInputSystem(inputChan chan PlayerInput) *PlayerInputSystem {
	return &PlayerInputSystem{InputChan: inputChan, Inputs: make(map[*websocket.Conn]PlayerInput), CloseSig: make(chan struct{})}
}

func (pis *PlayerInputSystem) Register(aether *ecs.Aether, hermes *ecs.Hermes, registrar *ecs.EntityRegistrar) {
	pis.aether = aether
	pis.hermes = hermes
	pis.registrar = registrar
}

func (pis *PlayerInputSystem) Init() {
	go pis.Collect()
}

func (pis *PlayerInputSystem) Update(dt float32) {
	inputs := pis.PopInputs()

	// Create mapping of Conn -> Entity
	stored, err := pis.aether.RetrieveType(reflect.TypeOf(&PlayerInput{}))
	if err != nil {
		panic(err)
	}
	connToEntity := make(map[*websocket.Conn]ecs.EntityId)
	for k, v := range stored {
		connToEntity[v.(*PlayerInput).Conn] = k
	}

	for conn, input := range inputs {
		entityId, exists := connToEntity[conn]
		if !exists {
			entityId = pis.registrar.NewId()
			pis.aether.Register(entityId, &input)
			pis.hermes.Send(&ecs.Message{Pipe: EntityCreatedPipe, EntityId: entityId})
		} else {
			pis.aether.Register(entityId, &input)
		}

		if input.Disconnected {
			close(input.OutputChan)
			pis.hermes.Send(&ecs.Message{Pipe: EntityDestroyedPipe, EntityId: entityId})
		}
	}
}

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
func (pis *PlayerInputSystem) ReceiveInput(input PlayerInput) {
	pis.mux.Lock()
	defer pis.mux.Unlock()
	pis.Inputs[input.Conn] = input
}

// PopInputs does a thread safe retrieve of player inputs since last Pop
func (pis *PlayerInputSystem) PopInputs() map[*websocket.Conn]PlayerInput {
	pis.mux.Lock()
	defer pis.mux.Unlock()

	inputs := pis.Inputs
	pis.Inputs = make(map[*websocket.Conn]PlayerInput)
	return inputs
}
