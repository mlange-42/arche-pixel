package plot_test

import (
	"testing"

	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/observer"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/plot"
	"github.com/mlange-42/arche-pixel/window"
	"github.com/stretchr/testify/assert"
)

func ExampleScatter() {

	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.TPS = 30

	// Create a scatter plot.
	m.AddUISystem((&window.Window{}).
		With(&plot.Scatter{
			Observers: []observer.Table{
				&TableObserver{}, // One or more observers.
			},
			X: []string{
				"X", // One X column per observer.
			},
			Y: [][]string{
				{"A", "B", "C"}, // One or more Y columns per observer.
			},
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

func TestScatter(t *testing.T) {
	m := model.New()
	m.TPS = 300

	m.AddUISystem((&window.Window{}).
		With(&plot.Scatter{
			Observers: []observer.Table{
				&TableObserver{},
			},
		}))

	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	m.Run()
}

func TestScatter_PanicXCount(t *testing.T) {
	m := model.New()
	m.TPS = 300

	m.AddUISystem((&window.Window{}).
		With(&plot.Scatter{
			Observers: []observer.Table{
				&TableObserver{},
			},
			X: []string{"X", "X"},
		}))

	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, m.Run)
}

func TestScatter_PanicX(t *testing.T) {
	m := model.New()
	m.TPS = 300

	m.AddUISystem((&window.Window{}).
		With(&plot.Scatter{
			Observers: []observer.Table{
				&TableObserver{},
			},
			X: []string{"F"},
		}))

	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, m.Run)
}

func TestScatter_PanicYCount(t *testing.T) {
	m := model.New()
	m.TPS = 300

	m.AddUISystem((&window.Window{}).
		With(&plot.Scatter{
			Observers: []observer.Table{
				&TableObserver{},
			},
			Y: [][]string{
				{"A"},
				{"A"},
			},
		}))

	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, m.Run)
}

func TestScatter_PanicY(t *testing.T) {
	m := model.New()
	m.TPS = 300

	m.AddUISystem((&window.Window{}).
		With(&plot.Scatter{
			Observers: []observer.Table{
				&TableObserver{},
			},
			Y: [][]string{
				{"F"},
			},
		}))

	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, m.Run)
}
