package plot_test

import (
	"testing"

	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/plot"
	"github.com/mlange-42/arche-pixel/window"
	"github.com/mlange-42/arche/ecs"
)

func ExampleResources() {
	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.TPS = 30

	// Create a window with a Resources drawer.
	m.AddUISystem((&window.Window{}).
		With(&plot.Resources{}))

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

func TestResources(t *testing.T) {
	m := model.New()

	_ = ecs.AddResource[Position](&m.World, &Position{})
	_ = ecs.ResourceID[Velocity](&m.World)

	m.TPS = 30

	m.AddUISystem((&window.Window{}).
		With(&plot.Resources{}))

	m.AddSystem(&system.FixedTermination{
		Steps: 10,
	})

	m.Run()
}
