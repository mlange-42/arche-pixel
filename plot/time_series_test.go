package plot_test

import (
	"math/rand"

	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/plot"
	"github.com/mlange-42/arche/ecs"
)

func ExampleTimeSeries() {

	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.Tps = 30
	m.Fps = 0

	// Create a time series plot.
	// See below for the implementation of the RowObserver.
	pl := plot.TimeSeries{
		Observer: &RowObserver{},
	}
	// Add the plot as UI system.
	m.AddUISystem(&pl)

	// Add a termination system that ends the simulation.
	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [github.com/faiface/pixel/pixelgl].
	// Uncomment the next line. It is commented out as the CI has no display device to test the model run.

	//pixelgl.Run(m.Run)
}

// RowObserver to generate random time series.
type RowObserver struct{}

func (o *RowObserver) Initialize(w *ecs.World) {}
func (o *RowObserver) Update(w *ecs.World)     {}
func (o *RowObserver) Header() []string {
	return []string{"A", "B", "C"}
}
func (o *RowObserver) Values(w *ecs.World) []float64 {
	return []float64{rand.Float64(), rand.Float64() + 1, rand.Float64() + 2}
}
