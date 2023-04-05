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
//
// See the example for [window] package.
type Drawer interface {
	Initialize(w *ecs.World, win *pixelgl.Window)
	Draw(w *ecs.World, win *pixelgl.Window)
}

// Window provides an OpenGL window for drawing.
//
// If the world contains a resource of type [github.com/mlange-42/arche-model/resource/Termination],
// the model is terminated when the window is closed.
//
// See the example for [window] package.
type Window struct {
	Bounds       Bounds
	DrawInterval int
	Drawers      []Drawer
	window       *pixelgl.Window
	step         int
	isClosed     bool
	termRes      generic.Resource[resource.Termination]
}

// AddDrawer adds a drawer.
func (s *Window) AddDrawer(d Drawer) {
	s.Drawers = append(s.Drawers, d)
}

// InitializeUI the system.
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
	s.step = 0
	s.isClosed = false
}

// UpdateUI the system.
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
	if s.DrawInterval <= 1 || s.step%s.DrawInterval == 0 {
		s.window.Clear(colornames.Black)

		for _, d := range s.Drawers {
			d.Draw(w, s.window)
		}
	}
	s.step++
}

// PostUpdateUI updates the GL window.
func (s *Window) PostUpdateUI(w *ecs.World) {
	s.window.Update()
}

// FinalizeUI the system.
func (s *Window) FinalizeUI(w *ecs.World) {}
