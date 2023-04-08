package plot_test

import (
	"math"

	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/observer"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/plot"
	"github.com/mlange-42/arche-pixel/window"
	"github.com/mlange-42/arche/ecs"
)

func ExampleImageRGB() {

	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.TPS = 30
	m.FPS = 0

	// Create an RGB image plot.
	// See below for the implementation of the CallbackMatrixObserver.
	m.AddUISystem((&window.Window{}).
		With(&plot.ImageRGB{
			Scale: 4,
			Observers: []observer.Matrix{
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return float64(i) / 240 }},
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return math.Sin(0.1 * float64(i)) }},
				&CallbackMatrixObserver{Callback: func(i, j int) float64 { return float64(j) / 160 }},
			},
			Min: []float64{0, 0, 0},
			Max: []float64{1, 1, 1},
		}))

	// Add a termination system that ends the simulation.
	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [github.com/faiface/pixel/pixelgl].
	// Uncomment the next line. It is commented out as the CI has no display device to test the model run.

	// pixelgl.Run(m.Run)
}

// Example observer, reporting a matrix filled with a callback(i, j).
type CallbackMatrixObserver struct {
	Callback func(i, j int) float64
	cols     int
	rows     int
	values   []float64
}

func (o *CallbackMatrixObserver) Initialize(w *ecs.World) {
	o.cols = 240
	o.rows = 160
	o.values = make([]float64, o.cols*o.rows)
}

func (o *CallbackMatrixObserver) Update(w *ecs.World) {}

func (o *CallbackMatrixObserver) Dims() (int, int) {
	return o.cols, o.rows
}

func (o *CallbackMatrixObserver) Values(w *ecs.World) []float64 {
	for idx := 0; idx < len(o.values); idx++ {
		i := idx % o.cols
		j := idx / o.cols
		o.values[idx] = o.Callback(i, j)
	}
	return o.values
}
