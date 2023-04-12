package window_test

import (
	"time"

	"github.com/faiface/pixel/pixelgl"
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/window"
)

func Example() {
	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.TPS = 30

	// Create a window system with a single drawer.
	win := (&window.Window{}).
		With(&RectDrawer{})
	// Add the window as UI system.
	m.AddUISystem(win)

	// Add a termination system that ends the simulation.
	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [github.com/faiface/pixel/pixelgl].
	// Note that the example will not work in the browser, as there is no proper display device available.
	pixelgl.Run(m.Run)
	time.Sleep(3 * time.Second)
	// Output:
}
