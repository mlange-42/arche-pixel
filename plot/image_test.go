package plot_test

import (
	"math"
	"time"

	"github.com/faiface/pixel/pixelgl"
	"github.com/mazznoer/colorgrad"
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/plot"
	"github.com/mlange-42/arche-pixel/window"
	"github.com/mlange-42/arche/ecs"
)

func ExampleImage() {

	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.TPS = 30
	m.FPS = 0

	// Create an image plot.
	// See below for the implementation of the MatrixObserver.
	m.AddUISystem(
		(&window.Window{}).
			With(&plot.Image{
				Scale:    4,
				Observer: &MatrixObserver{},
				Colors:   colorgrad.Inferno(),
				Min:      -2,
				Max:      2,
			}))

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

// Example observer, reporting a matrix with z = sin(0.1*i) + sin(0.2*j).
type MatrixObserver struct {
	cols   int
	rows   int
	values []float64
}

func (o *MatrixObserver) Initialize(w *ecs.World) {
	o.cols = 240
	o.rows = 160
	o.values = make([]float64, o.cols*o.rows)
}

func (o *MatrixObserver) Update(w *ecs.World) {}

func (o *MatrixObserver) Dims() (int, int) {
	return o.cols, o.rows
}

func (o *MatrixObserver) Values(w *ecs.World) []float64 {
	for idx := 0; idx < len(o.values); idx++ {
		i := idx % o.cols
		j := idx / o.cols
		o.values[idx] = math.Sin(0.1*float64(i)) + math.Sin(0.2*float64(j))
	}
	return o.values
}
