package ecs

import (
	"time"
)

// EntityId alias for entity ids
type EntityId uint64

// EntityManager manages the creation of new entities
type EntityRegistrar struct {
	LastId EntityId
}

// NewId generates a new unique id from an entity manager
func (registrar *EntityRegistrar) NewId() EntityId {
	registrar.LastId = registrar.LastId + 1
	return registrar.LastId
}

// System interface for types that run updates at every frame and can be ordered by priority
type System interface {
	Register(aether *Aether, hermes *Hermes, registrar *EntityRegistrar)
	Update(dt float64)
}

// Initializer interface for systems that can be initialized
type Initializer interface {
	Init()
}

// Cleaner interface for systems that run shutdowns
type Cleaner interface {
	Cleanup()
}

// GameLoop main game loop
func GameLoop(
	aether *Aether,
	hermes *Hermes,
	registrar *EntityRegistrar,
	systems []System,
	rate time.Duration,
	shutdown chan struct{}) {
	// Initialize server rate ticker
	ticker := time.NewTicker(rate)

	// Handle init for all
	for _, system := range systems {
		system.Register(aether, hermes, registrar)
		if initializer, ok := system.(Initializer); ok {
			initializer.Init()
		}
	}

	// Handle shutdown for all
	defer func() {
		ticker.Stop()
		for _, system := range systems {
			if cleaner, ok := system.(Cleaner); ok {
				cleaner.Cleanup()
			}
		}
		if shutdown != nil {
			close(shutdown)
		}
	}()

	// Game loop
	previous := time.Now()
	for {
		select {
		case <-shutdown:
			return
		case <-ticker.C:
			current := time.Now()
			dt := float64(current.Sub(previous)) / float64(time.Second)
			for _, system := range systems {
				system.Update(dt)
			}
			previous = current
		}
	}
}
