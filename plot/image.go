package plot

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mazznoer/colorgrad"
	"github.com/mlange-42/arche-model/observer"
	"github.com/mlange-42/arche/ecs"
)

// Image plot drawer.
type Image struct {
	Scale    float64            // Spatial scaling: cell size in screen pixels.
	Observer observer.Matrix    // Observer providing 2D matrix or grid data.
	Colors   colorgrad.Gradient // Colors for mapping values.
	Min      float64            // Minimum value for color mapping.
	Max      float64            // Maximum value for color mapping. Is set to 1.0 if both Min and Max are zero.
	slope    float64
	picture  *pixel.PictureData
}

// Initialize the system
func (s *Image) Initialize(w *ecs.World, win *pixelgl.Window) {
	s.Observer.Initialize(w)

	if s.Scale <= 0 {
		s.Scale = 1
	}
	if s.Min == 0 && s.Max == 0 {
		s.Max = 1
	}

	s.slope = 1.0 / (s.Max - s.Min)

	width, height := s.Observer.Dims()
	s.picture = pixel.MakePictureData(pixel.R(0, 0, float64(width), float64(height)))
}

// Update the drawer.
func (s *Image) Update(w *ecs.World) {
	s.Observer.Update(w)
}

// Draw the system
func (s *Image) Draw(w *ecs.World, win *pixelgl.Window) {
	values := s.Observer.Values(w)

	length := len(values)
	for i := 0; i < length; i++ {
		s.picture.Pix[i] = s.valueToColor(values[i])
	}

	sprite := pixel.NewSprite(s.picture, s.picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(s.picture.Rect.W()/2.0, s.picture.Rect.H()/2.0)).Scaled(pixel.Vec{}, s.Scale))
}

func (s *Image) valueToColor(v float64) color.RGBA {
	c := s.Colors.At((v - s.Min) * s.slope)
	return color.RGBA{
		R: uint8(c.R * 255),
		G: uint8(c.G * 255),
		B: uint8(c.B * 255),
		A: 0xff,
	}
}
