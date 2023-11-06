package plot_test

import (
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/plot"
	"github.com/mlange-42/arche-pixel/window"
)

func ExampleMonitor() {
	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.TPS = 30

	// Create a window with a Monitor drawer.
	m.AddUISystem((&window.Window{}).
		With(&plot.Monitor{}))

	// Add a termination system that ends the simulation.
	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [github.com/gopxl/pixel/v2/backends/opengl].
	// Uncomment the next line. It is commented out as the CI has no display device to test the model run.

	// opengl.Run(m.Run)

	// Output:
}

func ExampleNewMonitorWindow() {
	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.TPS = 30

	// Create a window with a Monitor drawer, using the shorthand constructor.
	m.AddUISystem(plot.NewMonitorWindow(10))

	// Add a termination system that ends the simulation.
	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [github.com/gopxl/pixel/v2/backends/opengl].
	// Uncomment the next line. It is commented out as the CI has no display device to test the model run.

	// opengl.Run(m.Run)

	// Output:
}
