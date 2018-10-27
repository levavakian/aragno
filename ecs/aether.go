package ecs

import (
	"errors"
	"reflect"
)

// Aether centralized storage for all components
type Aether struct {
	Maps map[reflect.Type]map[EntityId]interface{}
}

// NewAether creates a new aether
func NewAether() *Aether {
	return &Aether{Maps: make(map[reflect.Type]map[EntityId]interface{})}
}

// Register adds a component to the aether registry
func (ae *Aether) Register(id EntityId, cp interface{}) error {
	if reflect.ValueOf(cp).Kind() != reflect.Ptr {
		return errors.New("Aether.RegisterComponent: registered type must be pointer")
	}

	cpType := reflect.TypeOf(cp)

	if _, exists := ae.Maps[cpType]; !exists {
		ae.Maps[cpType] = make(map[EntityId]interface{})
	}

	ae.Maps[cpType][id] = cp
	return nil
}

// Retrieve retrieves a component based on type and entity id
func (ae *Aether) Retrieve(id EntityId, cpType reflect.Type) (interface{}, error) {
	if _, exists := ae.Maps[cpType]; !exists {
		if cpType.Kind() != reflect.Ptr {
			return nil, errors.New("Aether.Retrieve: not a pointer type")
		}
		return nil, errors.New("Aether.Retrieve: type not found")
	}

	if _, exists := ae.Maps[cpType][id]; !exists {
		return nil, errors.New("Aether.Retrieve: component not found for entity id")
	}

	return ae.Maps[cpType][id], nil
}

func (ae *Aether) RetrieveType(cpType reflect.Type) (map[EntityId]interface{}, error) {
	if _, exists := ae.Maps[cpType]; !exists {
		if cpType.Kind() != reflect.Ptr {
			return nil, errors.New("Aether.Retrieve: not a pointer type")
		}
		ae.Maps[cpType] = make(map[EntityId]interface{})
	}

	return ae.Maps[cpType], nil
}

// Deregister deregisters a component based on type and entity id
func (ae *Aether) DeregisterComponent(id EntityId, cpType reflect.Type) bool {
	if _, exists := ae.Maps[cpType]; !exists {
		return false
	}

	if _, exists := ae.Maps[cpType][id]; !exists {
		return false
	}

	delete(ae.Maps[cpType], id)
	return true
}

// DeregisterAll deregisters all components for a given id
func (ae *Aether) DeregisterAll(id EntityId) {
	for _, v := range ae.Maps {
		delete(v, id)
	}
}
