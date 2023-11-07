package plot_test

import (
	"math/rand"

	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/plot"
	"github.com/mlange-42/arche-pixel/window"
	"github.com/mlange-42/arche/ecs"
)

func ExampleLines() {

	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.TPS = 30

	// Create a line plot.
	// See below for the implementation of the TableObserver.
	m.AddUISystem((&window.Window{}).
		With(&plot.Lines{
			Observer: &TableObserver{},
			X:        "X",                     // Optional, defaults to row index
			Y:        []string{"A", "B", "C"}, // Optional, defaults to all but X
		}))

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

// TableObserver to generate random time series.
type TableObserver struct {
	data [][]float64
}

func (o *TableObserver) Initialize(w *ecs.World) {
	rows := 25
	o.data = make([][]float64, rows)

	for i := 0; i < rows; i++ {
		o.data[i] = []float64{float64(i), float64(i) / float64(rows), float64(rows-i) / float64(rows), 0}
	}
}
func (o *TableObserver) Update(w *ecs.World) {}
func (o *TableObserver) Header() []string {
	return []string{"X", "A", "B", "C"}
}
func (o *TableObserver) Values(w *ecs.World) [][]float64 {
	for i := 0; i < len(o.data); i++ {
		o.data[i][3] = rand.Float64()
	}
	return o.data
}
