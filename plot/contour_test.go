package plot_test

import (
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/observer"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/plot"
	"github.com/mlange-42/arche-pixel/window"
	"gonum.org/v1/plot/palette"
)

func ExampleContour() {
	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.TPS = 30
	m.FPS = 0

	// Create a contour plot.
	m.AddUISystem(
		(&window.Window{}).
			With(&plot.Contour{
				Observer: observer.MatrixToGrid(&MatrixObserver{}, nil, nil),
				Palette:  palette.Heat(16, 1),
				Levels:   []float64{-2, -1.5, -1, -0.5, 0, 0.5, 1, 1.5, 2},
			}))

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
