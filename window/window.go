package window

import (
	"fmt"
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mlange-42/arche-model/resource"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"golang.org/x/image/colornames"
)

// Drawer interface.
// Drawers are used by the [Window] to render information from an Arche model.
type Drawer interface {
	Initialize(w *ecs.World, win *pixelgl.Window) // Initialize the Drawer.
	Update(w *ecs.World)                          // Update is called by normal system update.
	Draw(w *ecs.World, win *pixelgl.Window)       // Draw is called by the UI systems update.
}

// Bounds define a bounding box for a window.
type Bounds struct {
	X      int
	Y      int
	Width  int
	Height int
}

// B created a new Bounds object.
func B(x, y, w, h int) Bounds {
	return Bounds{x, y, w, h}
}

// Window provides an OpenGL window for drawing.
// Drawing is done by one or more [Drawer] instances.
// Further, window bounds and update and draw intervals van be configured.
//
// If the world contains a resource of type [github.com/mlange-42/arche-model/resource/Termination],
// the model is terminated when the window is closed.
type Window struct {
	Title        string   // Window title. Optional.
	Bounds       Bounds   // Window bounds (position and size). Optional.
	Drawers      []Drawer // Drawers in increasing z order.
	DrawInterval int      // Interval for re-drawing, in UI frames. Optional.
	window       *pixelgl.Window
	drawStep     int64
	isClosed     bool
	termRes      generic.Resource[resource.Termination]
}

// With adds one or more [Drawer] instances to the window.
func (w *Window) With(drawers ...Drawer) *Window {
	w.Drawers = append(w.Drawers, drawers...)
	return w
}

// Initialize the window system.
func (w *Window) Initialize(world *ecs.World) {}

// InitializeUI the window system.
func (w *Window) InitializeUI(world *ecs.World) {
	if w.Bounds.Width <= 0 {
		w.Bounds.Width = 1024
	}
	if w.Bounds.Height <= 0 {
		w.Bounds.Height = 768
	}
	if w.Title == "" {
		w.Title = "Arche"
	}
	cfg := pixelgl.WindowConfig{
		Title:     w.Title,
		Bounds:    pixel.R(0, 0, float64(w.Bounds.Width), float64(w.Bounds.Height)),
		Position:  pixel.V(float64(w.Bounds.X), float64(w.Bounds.Y)),
		Resizable: true,
	}

	defer func() {
		if err := recover(); err != nil {
			txt := fmt.Sprint(err)
			if txt == "mainthread: did not call Run" {
				log.Fatal("ERROR: when using graphics via the pixel engine, run the model like this:\n    pixelgl.Run(model.Run)")
			}
			panic(err)
		}
	}()

	var err error
	w.window, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	for _, d := range w.Drawers {
		d.Initialize(world, w.window)
	}

	w.termRes = generic.NewResource[resource.Termination](world)
	w.drawStep = 0
	w.isClosed = false
}

// Update the window system.
func (w *Window) Update(world *ecs.World) {
	if w.isClosed {
		return
	}
	for _, d := range w.Drawers {
		d.Update(world)
	}
}

// UpdateUI the window system.
func (w *Window) UpdateUI(world *ecs.World) {
	if w.window.Closed() {
		if !w.isClosed {
			term := w.termRes.Get()
			if term != nil {
				term.Terminate = true
			}
			w.isClosed = true
		}
		return
	}
	if w.DrawInterval <= 1 || w.drawStep%int64(w.DrawInterval) == 0 {
		w.window.Clear(colornames.Black)

		for _, d := range w.Drawers {
			d.Draw(world, w.window)
		}
	}
	w.drawStep++
}

// PostUpdateUI updates the underlying GL window.
func (w *Window) PostUpdateUI(world *ecs.World) {
	w.window.Update()
}

// Finalize the window system.
func (w *Window) Finalize(world *ecs.World) {}

// FinalizeUI the window system.
func (w *Window) FinalizeUI(world *ecs.World) {}
