package plot

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mlange-42/arche-model/observer"
	"github.com/mlange-42/arche/ecs"
)

// ImageRGB plot drawer.
type ImageRGB struct {
	Scale     float64           // Spatial scaling: cell size in screen pixels.
	Observers []observer.Matrix // Observers for the red, green and blue channel. Elements can be nil.
	Min       []float64         // Minimum value for channel color mapping. Default [0, 0, 0].
	Max       []float64         // Maximum value for channel color mapping. Default [1, 1, 1].
	slope     []float64
	dataLen   int
	picture   *pixel.PictureData
}

// Initialize the drawer.
func (s *ImageRGB) Initialize(w *ecs.World, win *pixelgl.Window) {
	if len(s.Observers) != 3 {
		panic("RgbImage plot needs exactly 3 observers")
	}
	width, height := -1, -1
	for i := 0; i < 3; i++ {
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

	s.slope = []float64{
		1.0 / (s.Max[0] - s.Min[0]),
		1.0 / (s.Max[1] - s.Min[1]),
		1.0 / (s.Max[2] - s.Min[2]),
	}

	s.dataLen = width * height
	s.picture = pixel.MakePictureData(pixel.R(0, 0, float64(width), float64(height)))
}

// Update the drawer.
func (s *ImageRGB) Update(w *ecs.World) {
	for i := 0; i < 3; i++ {
		if s.Observers[i] != nil {
			s.Observers[i].Update(w)
		}
	}
}

// Draw the drawer.
func (s *ImageRGB) Draw(w *ecs.World, win *pixelgl.Window) {
	cannels := make([][]float64, 3)
	for i := 0; i < 3; i++ {
		if s.Observers[i] != nil {
			cannels[i] = s.Observers[i].Values(w)
		}
	}

	values := append([]float64{}, s.Min...)
	for i := 0; i < s.dataLen; i++ {
		for j := 0; j < 3; j++ {
			if cannels[j] != nil {
				values[j] = cannels[j][i]
			}
		}
		s.picture.Pix[i] = s.valuesToColor(values[0], values[1], values[2])
	}

	sprite := pixel.NewSprite(s.picture, s.picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(s.picture.Rect.W()/2.0, s.picture.Rect.H()/2.0)).Scaled(pixel.Vec{}, s.Scale))
}

func (s *ImageRGB) valuesToColor(r, g, b float64) color.RGBA {
	return color.RGBA{
		R: norm(r, s.Min[0], s.slope[0]),
		G: norm(g, s.Min[1], s.slope[1]),
		B: norm(b, s.Min[2], s.slope[2]),
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
