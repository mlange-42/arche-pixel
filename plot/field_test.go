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

func ExampleField() {
	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.TPS = 30
	m.FPS = 0

	// Create a contour plot.
	m.AddUISystem(
		(&window.Window{}).
			With(&plot.Field{
				Observer: observer.LayersToLayers(&FieldObserver{}, nil, nil),
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

type FieldObserver struct {
	cols   int
	rows   int
	values [][]float64
}

func (o *FieldObserver) Initialize(w *ecs.World) {
	o.cols = 60
	o.rows = 40
	o.values = make([][]float64, 2)
	for i := 0; i < len(o.values); i++ {
		o.values[i] = make([]float64, o.cols*o.rows)
	}
}

func (o *FieldObserver) Update(w *ecs.World) {}

func (o *FieldObserver) Dims() (int, int) {
	return o.cols, o.rows
}

func (o *FieldObserver) Layers() int {
	return 2
}

func (o *FieldObserver) Values(w *ecs.World) [][]float64 {
	ln := len(o.values[0])
	for idx := 0; idx < ln; idx++ {
		i := idx % o.cols
		j := idx / o.cols
		o.values[0][idx] = math.Sin(float64(i))
		o.values[1][idx] = -math.Sin(float64(j))
	}
	return o.values
}
