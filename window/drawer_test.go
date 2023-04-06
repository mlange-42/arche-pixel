package window_test

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mlange-42/arche-model/resource"
	"github.com/mlange-42/arche-pixel/window"
	"github.com/mlange-42/arche/ecs"
)

func ExampleDrawer() {
	var dr window.Drawer = &RectDrawer{}
	_ = dr
	// Output:
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

// Update the RectDrawer (does nothing).
func (d *RectDrawer) Update(w *ecs.World) {}

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
