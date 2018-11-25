package game

import (
	"aragno/ecs"
	"aragno/game/component"
	"aragno/game/system"
	"time"
)

// Loop registers systems and game world, and loops until shutdown
func Loop(input chan component.PlayerInput) {
	systems := make([]ecs.System, 0)
	rate := time.Millisecond * 16 // 60Hz

	systems = append(systems,
		system.NewPlayerInputSystem(input),
		system.NewMapSystem(),
		//		system.NewMoveSystem(),
		system.NewDynamicsSystem(),
		system.NewStateOutputSystem(),
		system.NewEntityDestroyerSystem())

	aether := ecs.NewAether(true)
	hermes := ecs.NewHermes()
	registrar := &ecs.EntityRegistrar{}

	ecs.GameLoop(aether, hermes, registrar, systems, rate, nil)
}
