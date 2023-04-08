package window_test

import (
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-pixel/window"
)

func ExampleWindow() {
	m := model.New()

	// Create a Window with at least one Drawer.
	window := (&window.Window{Bounds: window.B(100, 100, 800, 600)}).
		With(&RectDrawer{})

	// Add is to the model as UI system.
	m.AddUISystem(window)
	// Output:
}
