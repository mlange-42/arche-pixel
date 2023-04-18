package plot

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mlange-42/arche-model/observer"
	"github.com/mlange-42/arche-pixel/window"
	"github.com/mlange-42/arche/ecs"
)

// ImageRGB drawer.
//
// Draws an image from a Matrix observer per RGB color channel.
// The image is scaled to the canvas extent, with preserved aspect ratio.
// Does not add plot axes etc.
type ImageRGB struct {
	Scale     float64           // Spatial scaling: cell size in screen pixels. Optional, default auto.
	Observers []observer.Matrix // Observers for the red, green and blue channel. Elements can be nil.
	Min       []float64         // Minimum value for channel color mapping. Optional, default [0, 0, 0].
	Max       []float64         // Maximum value for channel color mapping. Optional, default [1, 1, 1].
	slope     []float64
	dataLen   int
	picture   *pixel.PictureData
}

// Initialize the drawer.
func (i *ImageRGB) Initialize(w *ecs.World, win *pixelgl.Window) {
	if len(i.Observers) != 3 {
		panic("RgbImage plot needs exactly 3 observers")
	}
	width, height := -1, -1
	for j := 0; j < 3; j++ {
		if i.Observers[j] == nil {
			continue
		}
		i.Observers[j].Initialize(w)
		wi, he := i.Observers[j].Dims()
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

	if i.Min == nil {
		i.Min = []float64{0, 0, 0}
	}
	if i.Max == nil {
		i.Max = []float64{1, 1, 1}
	}
	if len(i.Min) != 3 {
		panic("RgbImage plot needs exactly 3 Min values")
	}
	if len(i.Max) != 3 {
		panic("RgbImage plot needs exactly 3 Max values")
	}

	i.slope = []float64{
		1.0 / (i.Max[0] - i.Min[0]),
		1.0 / (i.Max[1] - i.Min[1]),
		1.0 / (i.Max[2] - i.Min[2]),
	}

	i.dataLen = width * height
	i.picture = pixel.MakePictureData(pixel.R(0, 0, float64(width), float64(height)))
}

// Update the drawer.
func (i *ImageRGB) Update(w *ecs.World) {
	for j := 0; j < 3; j++ {
		if i.Observers[j] != nil {
			i.Observers[j].Update(w)
		}
	}
}

// UpdateInputs handles input events of the previous frame update.
func (i *ImageRGB) UpdateInputs(w *ecs.World, win *pixelgl.Window) {}

// Draw the drawer.
func (i *ImageRGB) Draw(w *ecs.World, win *pixelgl.Window) {
	cannels := make([][]float64, 3)
	for j := 0; j < 3; j++ {
		if i.Observers[j] != nil {
			cannels[j] = i.Observers[j].Values(w)
		}
	}

	values := append([]float64{}, i.Min...)
	for j := 0; j < i.dataLen; j++ {
		for k := 0; k < 3; k++ {
			if cannels[k] != nil {
				values[k] = cannels[k][j]
			}
		}
		i.picture.Pix[j] = i.valuesToColor(values[0], values[1], values[2])
	}

	scale := i.Scale
	if i.Scale <= 0 {
		scale = window.Scale(win, i.picture.Rect.W(), i.picture.Rect.H())
	}

	sprite := pixel.NewSprite(i.picture, i.picture.Bounds())
	sprite.Draw(win,
		pixel.IM.Moved(pixel.V(i.picture.Rect.W()/2.0, i.picture.Rect.H()/2.0)).
			Scaled(pixel.Vec{}, scale),
	)
}

func (i *ImageRGB) valuesToColor(r, g, b float64) color.RGBA {
	return color.RGBA{
		R: norm(r, i.Min[0], i.slope[0]),
		G: norm(g, i.Min[1], i.slope[1]),
		B: norm(b, i.Min[2], i.slope[2]),
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
