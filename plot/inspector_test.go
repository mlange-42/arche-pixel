package plot_test

import (
	"testing"

	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/resource"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/plot"
	"github.com/mlange-42/arche-pixel/window"
	"github.com/mlange-42/arche/ecs"
)

func ExampleInspector() {
	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.TPS = 30

	// Create an entity to inspect it.
	posID := ecs.ComponentID[Position](&m.World)
	velID := ecs.ComponentID[Velocity](&m.World)
	entity := m.World.NewEntity(posID, velID)

	// Set it as selected entity.
	ecs.AddResource(&m.World, &resource.SelectedEntity{Selected: entity})

	// Create a window with an Inspector drawer.
	m.AddUISystem((&window.Window{}).
		With(&plot.Inspector{}))

	// Add a termination system that ends the simulation.
	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	m.Run()

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [window.Run].
	// Comment out the code line above, and uncomment the next line to run this example stand-alone.

	// window.Run(m)

	// Output:
}

func TestInspector(t *testing.T) {
	m := model.New()
	m.TPS = 300

	posID := ecs.ComponentID[Position](&m.World)
	velID := ecs.ComponentID[Velocity](&m.World)
	entity := m.World.NewEntity(posID, velID)

	ecs.AddResource(&m.World, &resource.SelectedEntity{Selected: entity})

	m.AddUISystem((&window.Window{}).
		With(&plot.Inspector{
			HideNames: true,
		}))

	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	m.Run()
}

func TestInspector_DeadEntity(t *testing.T) {
	m := model.New()
	m.TPS = 300

	posID := ecs.ComponentID[Position](&m.World)
	velID := ecs.ComponentID[Velocity](&m.World)
	entity := m.World.NewEntity(posID, velID)

	ecs.AddResource(&m.World, &resource.SelectedEntity{Selected: entity})

	m.AddUISystem((&window.Window{}).
		With(&plot.Inspector{
			HideNames: true,
		}))

	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	m.World.RemoveEntity(entity)

	m.Run()
}
