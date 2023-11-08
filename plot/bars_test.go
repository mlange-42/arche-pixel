package plot_test

import (
	"testing"

	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/plot"
	"github.com/mlange-42/arche-pixel/window"
	"github.com/stretchr/testify/assert"
)

func ExampleBars() {

	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.TPS = 30

	// Create a time series plot.
	m.AddUISystem((&window.Window{}).
		With(&plot.Bars{
			Observer: &RowObserver{},
			YLim:     [...]float64{0, 4}, // Optional Y axis limits.
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

func TestBars_Columns(t *testing.T) {
	m := model.New()
	m.TPS = 300
	m.AddUISystem((&window.Window{}).
		With(&plot.Bars{
			Observer: &RowObserver{},
			YLim:     [...]float64{0, 4},
			Columns:  []string{"A", "C"},
		}))
	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	m.Run()
}

func TestBars_PanicColumns(t *testing.T) {
	m := model.New()
	m.TPS = 300
	m.AddUISystem((&window.Window{}).
		With(&plot.Bars{
			Observer: &RowObserver{},
			Columns:  []string{"A", "F"},
		}))
	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})
	assert.Panics(t, m.Run)
}
