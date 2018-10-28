package game

import (
	"aragno/ecs"
	"time"
)

func Loop(input chan PlayerInput) {
	systems := make([]ecs.System, 0)
	rate := time.Millisecond * 16 // 60Hz

	systems = append(systems,
		NewPlayerInputSystem(input),
		NewStateOutputSystem(),
		NewEntityDestroyerSystem())

	ecs.GameLoop(systems, rate, nil)
}
