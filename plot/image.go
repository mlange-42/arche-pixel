package plot

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mazznoer/colorgrad"
	"github.com/mlange-42/arche-model/observer"
	"github.com/mlange-42/arche-pixel/window"
	"github.com/mlange-42/arche/ecs"
)

// Image plot reporter.
//
// If the world contains a resource of type [github.com/mlange-42/arche-model/resource.Termination],
// the model is terminated when the window is closed.
type Image struct {
	Bounds         window.Bounds      // Window bounds.
	Scale          float64            // Spatial scaling: cell size in screen pixels.
	Observer       observer.Matrix    // The observer.
	Colors         colorgrad.Gradient // Colors for mapping values.
	Min            float64            // Minimum value for color mapping.
	Max            float64            // Maximum value for color mapping. Is set to 1.0 if both Min and Max are zero.
	UpdateInterval int                // Interval for updating the observer, in model ticks.
	DrawInterval   int                // Interval for re-drawing, in UI frames.
	window.Window
	drawer imageDrawer
	step   int64
}

// Initialize the system.
func (s *Image) Initialize(w *ecs.World) {
	s.step = 0
}

// InitializeUI the system.
func (s *Image) InitializeUI(w *ecs.World) {
	s.Observer.Initialize(w)

	if s.Scale <= 0 {
		s.Scale = 1
	}
	if s.Min == 0 && s.Max == 0 {
		s.Max = 1
	}

	s.drawer = newImageDrawer(s.Observer, s.Scale, &s.Colors, s.Min, s.Max)

	s.Window.DrawInterval = s.DrawInterval
	s.Window.Bounds = s.Bounds
	s.Window.Drawers = append([]window.Drawer{&s.drawer}, s.Window.Drawers...)
	s.Window.InitializeUI(w)
}

// Update the system.
func (s *Image) Update(w *ecs.World) {
	if s.UpdateInterval <= 1 || s.step%int64(s.UpdateInterval) == 0 {
		s.Observer.Update(w)
	}
	s.step++
}

// Finalize the system.
func (s *Image) Finalize(w *ecs.World) {}

type imageDrawer struct {
	observer observer.Matrix
	scale    float64
	colors   *colorgrad.Gradient
	offset   float64
	slope    float64
	picture  *pixel.PictureData
}

func newImageDrawer(obs observer.Matrix, scale float64, colors *colorgrad.Gradient, min, max float64) imageDrawer {
	return imageDrawer{
		observer: obs,
		scale:    scale,
		colors:   colors,
		offset:   min,
		slope:    1.0 / (max - min),
	}
}

// Initialize the system
func (s *imageDrawer) Initialize(w *ecs.World, win *pixelgl.Window) {
	width, height := s.observer.Dims()
	s.picture = pixel.MakePictureData(pixel.R(0, 0, float64(width), float64(height)))
}

// Draw the system
func (s *imageDrawer) Draw(w *ecs.World, win *pixelgl.Window) {
	values := s.observer.Values(w)

	length := len(values)
	for i := 0; i < length; i++ {
		s.picture.Pix[i] = s.valueToColor(values[i])
	}

	sprite := pixel.NewSprite(s.picture, s.picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(s.picture.Rect.W()/2.0, s.picture.Rect.H()/2.0)).Scaled(pixel.Vec{}, s.scale))
}

func (s *imageDrawer) valueToColor(v float64) color.RGBA {
	c := s.colors.At((v - s.offset) * s.slope)
	return color.RGBA{
		R: uint8(c.R * 255),
		G: uint8(c.G * 255),
		B: uint8(c.B * 255),
		A: 0xff,
	}
}
