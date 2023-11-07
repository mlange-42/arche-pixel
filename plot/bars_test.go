package plot_test

import (
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/plot"
	"github.com/mlange-42/arche-pixel/window"
)

func ExampleBars() {

	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.TPS = 30

	// Create a time series plot.
	m.AddUISystem((&window.Window{}).
		With(&plot.Bars{
			Observer: &RowObserver{},
			YLim:     [...]float64{0, 4}, // Optional Y axis limits.
		}))

	// Add a termination system that ends the simulation.
	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [window.Run].
	// Uncomment the next line. It is commented out as the CI has no display device to test the model run.

	window.Run(m)

	// Output:
}
