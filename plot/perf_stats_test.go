package plot_test

import (
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/plot"
	"github.com/mlange-42/arche-pixel/window"
)

func ExamplePerfStats() {
	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.TPS = 30

	// Create a window with a PerfStats drawer.
	// PerfStats is intended as an overlay, so more drawers can be added before it.
	m.AddUISystem((&window.Window{}).
		With(&plot.PerfStats{}))

	// Add a termination system that ends the simulation.
	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [window.Run].
	// Uncomment the next line. It is commented out as the CI has no display device to test the model run.

	// window.Run(m)

	// Output:
}
