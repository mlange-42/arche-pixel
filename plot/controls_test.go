package plot_test

import (
	"time"

	"github.com/faiface/pixel/pixelgl"
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/plot"
	"github.com/mlange-42/arche-pixel/window"
)

func ExampleControls() {
	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.TPS = 30

	// Create a window with a Controls drawer.
	m.AddUISystem((&window.Window{}).
		With(&plot.Controls{Scale: 2}))

	// Alternatively, add the Controls in addition to another drawer (e.g. Monitor).
	m.AddUISystem((&window.Window{}).
		With(
			&plot.Monitor{},
			&plot.Controls{},
		))

	// Add a termination system that ends the simulation.
	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [github.com/faiface/pixel/pixelgl].
	// Note that the example will not work in the browser, as there is no proper display device available.
	pixelgl.Run(m.Run)
	time.Sleep(10 * time.Second)
	// Output:
}
