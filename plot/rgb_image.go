package plot

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mlange-42/arche-model/observer"
	"github.com/mlange-42/arche-pixel/window"
	"github.com/mlange-42/arche/ecs"
)

// ImageRGB plot reporter.
//
// If the world contains a resource of type [github.com/mlange-42/arche-model/resource.Termination],
// the model is terminated when the window is closed.
type ImageRGB struct {
	Bounds         window.Bounds     // Window bounds.
	Scale          float64           // Spatial scaling: cell size in screen pixels.
	Observers      []observer.Matrix // Observers for the red, green and blue channel. Elements can be nil.
	Min            []float64         // Minimum value for channel color mapping. Default [0, 0, 0].
	Max            []float64         // Maximum value for channel color mapping. Default [1, 1, 1].
	UpdateInterval int               // Interval for updating the observer, in model ticks.
	DrawInterval   int               // Interval for re-drawing, in UI frames.
	window.Window
	drawer rgbImageDrawer
	step   int64
}

// Initialize the system.
func (s *ImageRGB) Initialize(w *ecs.World) {
	s.step = 0
}

// InitializeUI the system.
func (s *ImageRGB) InitializeUI(w *ecs.World) {
	if len(s.Observers) != 3 {
		panic("RgbImage plot needs exactly 3 observers")
	}
	width, height := -1, -1
	for i := 0; i < len(s.Observers); i++ {
		if s.Observers[i] == nil {
			continue
		}
		s.Observers[i].Initialize(w)
		wi, he := s.Observers[i].Dims()
		if width >= 0 && width != wi {
			panic("observers differ in matrix width")
		}
		if height >= 0 && height != he {
			panic("observers differ in matrix width")
		}
		width, height = wi, he
	}
	if width < 0 && height < 0 {
		panic("needs an observer for at least one channel")
	}

	if s.Scale <= 0 {
		s.Scale = 1
	}
	if s.Min == nil {
		s.Min = []float64{0, 0, 0}
	}
	if s.Max == nil {
		s.Max = []float64{1, 1, 1}
	}
	if len(s.Min) != 3 {
		panic("RgbImage plot needs exactly 3 Min values")
	}
	if len(s.Max) != 3 {
		panic("RgbImage plot needs exactly 3 Max values")
	}

	s.drawer = newRgbImageDrawer(s.Observers, s.Scale, s.Min, s.Max)

	s.Window.DrawInterval = s.DrawInterval
	s.Window.Bounds = s.Bounds
	s.Window.Drawers = append([]window.Drawer{&s.drawer}, s.Window.Drawers...)
	s.Window.InitializeUI(w)
}

// Update the system.
func (s *ImageRGB) Update(w *ecs.World) {
	if s.UpdateInterval <= 1 || s.step%int64(s.UpdateInterval) == 0 {
		for i := 0; i < len(s.Observers); i++ {
			if s.Observers[i] != nil {
				s.Observers[i].Update(w)
			}
		}
	}
	s.step++
}

// Finalize the system.
func (s *ImageRGB) Finalize(w *ecs.World) {}

type rgbImageDrawer struct {
	observers []observer.Matrix
	offset    []float64
	slope     []float64
	scale     float64
	dataLen   int
	picture   *pixel.PictureData
}

func newRgbImageDrawer(obs []observer.Matrix, scale float64, min, max []float64) rgbImageDrawer {
	slope := []float64{
		1.0 / (max[0] - min[0]),
		1.0 / (max[1] - min[1]),
		1.0 / (max[2] - min[2]),
	}
	return rgbImageDrawer{
		observers: obs,
		scale:     scale,
		offset:    min,
		slope:     slope,
	}
}

// Initialize the system
func (s *rgbImageDrawer) Initialize(w *ecs.World, win *pixelgl.Window) {
	width, height := -1, -1
	for i := 0; i < 3; i++ {
		if s.observers[i] != nil {
			width, height = s.observers[i].Dims()
			break
		}
	}
	if width < 0 && height < 0 {
		panic("needs an observer for at least one channel")
	}
	s.dataLen = width * height
	s.picture = pixel.MakePictureData(pixel.R(0, 0, float64(width), float64(height)))
}

// Draw the system
func (s *rgbImageDrawer) Draw(w *ecs.World, win *pixelgl.Window) {
	cannels := make([][]float64, 3)
	for i := 0; i < 3; i++ {
		if s.observers[i] != nil {
			cannels[i] = s.observers[i].Values(w)
		}
	}

	values := append([]float64{}, s.offset...)
	for i := 0; i < s.dataLen; i++ {
		for j := 0; j < 3; j++ {
			if cannels[j] != nil {
				values[j] = cannels[j][i]
			}
		}
		s.picture.Pix[i] = s.valuesToColor(values[0], values[1], values[2])
	}

	sprite := pixel.NewSprite(s.picture, s.picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(s.picture.Rect.W()/2.0, s.picture.Rect.H()/2.0)).Scaled(pixel.Vec{}, s.scale))
}

func (s *rgbImageDrawer) valuesToColor(r, g, b float64) color.RGBA {
	return color.RGBA{
		R: norm(r, s.offset[0], s.slope[0]),
		G: norm(g, s.offset[1], s.slope[1]),
		B: norm(b, s.offset[2], s.slope[2]),
		A: 0xff,
	}
}

func norm(v, off, slope float64) uint8 {
	vv := (v - off) * slope
	if vv <= 0 {
		return 0
	}
	if vv >= 1 {
		return 255
	}
	return uint8(vv * 255)
}
