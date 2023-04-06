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

// Drawer interface.
// Drawers are used by the [Window] to render information from an Arche model.
type Drawer interface {
	Initialize(w *ecs.World, win *pixelgl.Window) // Initialize the Drawer.
	Update(w *ecs.World)                          // Update is called by normal system update.
	Draw(w *ecs.World, win *pixelgl.Window)       // Draw is called by the UI systems update.
}

// Window provides an OpenGL window for drawing.
// Drawing is done by one or more [Drawer] instances.
// Further, window bounds and update and draw intervals van be configured.
//
// If the world contains a resource of type [github.com/mlange-42/arche-model/resource/Termination],
// the model is terminated when the window is closed.
type Window struct {
	Bounds         Bounds   // Window bounds (position and size). Optional.
	Drawers        []Drawer // Drawers in increasing z order.
	UpdateInterval int      // Interval for updating drawers (and thus potentially observers), in model ticks. Optional.
	DrawInterval   int      // Interval for re-drawing, in UI frames. Optional.
	window         *pixelgl.Window
	updateStep     int64
	drawStep       int64
	isClosed       bool
	termRes        generic.Resource[resource.Termination]
}

// AddDrawer adds a [Drawer] to the window.
func (s *Window) AddDrawer(d Drawer) {
	s.Drawers = append(s.Drawers, d)
}

// Initialize the window system.
func (s *Window) Initialize(w *ecs.World) {
	s.updateStep = 0
}

// InitializeUI the window system.
func (s *Window) InitializeUI(w *ecs.World) {
	if s.Bounds.Width <= 0 {
		s.Bounds.Width = 1024
	}
	if s.Bounds.Height <= 0 {
		s.Bounds.Height = 768
	}
	cfg := pixelgl.WindowConfig{
		Title:    "Arche",
		Bounds:   pixel.R(0, 0, float64(s.Bounds.Width), float64(s.Bounds.Height)),
		Position: pixel.V(float64(s.Bounds.X), float64(s.Bounds.Y)),
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
	s.window, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	for _, d := range s.Drawers {
		d.Initialize(w, s.window)
	}

	s.termRes = generic.NewResource[resource.Termination](w)
	s.drawStep = 0
	s.isClosed = false
}

// Update the window system.
func (s *Window) Update(w *ecs.World) {
	if s.isClosed {
		return
	}
	if s.UpdateInterval <= 1 || s.updateStep%int64(s.UpdateInterval) == 0 {
		for _, d := range s.Drawers {
			d.Update(w)
		}
	}
	s.updateStep++
}

// UpdateUI the window system.
func (s *Window) UpdateUI(w *ecs.World) {
	if s.window.Closed() {
		if !s.isClosed {
			term := s.termRes.Get()
			if term != nil {
				term.Terminate = true
			}
			s.isClosed = true
		}
		return
	}
	if s.DrawInterval <= 1 || s.drawStep%int64(s.DrawInterval) == 0 {
		s.window.Clear(colornames.Black)

		for _, d := range s.Drawers {
			d.Draw(w, s.window)
		}
	}
	s.drawStep++
}

// PostUpdateUI updates the underlying GL window.
func (s *Window) PostUpdateUI(w *ecs.World) {
	s.window.Update()
}

// Finalize the window system.
func (s *Window) Finalize(w *ecs.World) {}

// FinalizeUI the window system.
func (s *Window) FinalizeUI(w *ecs.World) {}
