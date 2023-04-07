package plot_test

import (
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/plot"
	"github.com/mlange-42/arche-pixel/window"
)

func ExampleWorldStats() {
	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.Tps = 30
	m.Fps = 0

	// Create an image plot.
	// See below for the implementation of the MatrixObserver.
	m.AddUISystem(&window.Window{
		Drawers: []window.Drawer{
			&plot.WorldStats{},
		},
	})

	// Add a termination system that ends the simulation.
	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [github.com/faiface/pixel/pixelgl].
	// Uncomment the next line. It is commented out as the CI has no display device to test the model run.

	// pixelgl.Run(m.Run)
}
