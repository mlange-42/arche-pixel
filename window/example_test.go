package window_test

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/resource"
	"github.com/mlange-42/arche-model/system"
	"github.com/mlange-42/arche-pixel/window"
	"github.com/mlange-42/arche/ecs"
)

func Example() {
	// Create a new model.
	m := model.New()

	// Limit the the simulation speed.
	m.Tps = 30
	m.Fps = 0

	// Create a window system with a single drawer.
	win := window.Window{
		Drawers: []window.Drawer{&RectDrawer{}},
	}
	// Add the window as UI system.
	m.AddUISystem(&win)

	// Add a termination system that ends the simulation.
	m.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	// Due to the use of the OpenGL UI system, the model must be run via [github.com/faiface/pixel/pixelgl].
	// Uncomment the next line. It is commented out as the CI has no display device to test the model run.

	//pixelgl.Run(m.Run)
}

// RectDrawer is an example drawer.
type RectDrawer struct {
	dr imdraw.IMDraw
}

// Initialize the RectDrawer.
func (d *RectDrawer) Initialize(w *ecs.World, win *pixelgl.Window) {
	// Create a drawer from the Pixel engine.
	d.dr = *imdraw.New(nil)
}

// Draw the RectDrawer's stuff.
func (d *RectDrawer) Draw(w *ecs.World, win *pixelgl.Window) {
	// Get a resource from the world.
	tick := ecs.GetResource[resource.Tick](w)
	offset := float64(tick.Tick)

	// Create a white rectangle that moves with progressing model tick.
	d.dr.Color = color.White
	d.dr.Push(pixel.V(50+offset, 50+offset), pixel.V(250+offset, 200+offset))
	d.dr.Rectangle(0)

	// Draw everything on the window.
	d.dr.Draw(win)

	// Reset the drawer
	d.dr.Reset()
	d.dr.Clear()
}
