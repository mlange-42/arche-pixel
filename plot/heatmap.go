package plot

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mlange-42/arche-model/observer"
	"github.com/mlange-42/arche/ecs"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

// HeatMap plot drawer.
//
// Plots a grid as a heatmap image.
// For large grids, this is relatively slow.
// Consider using [Image] instead.
type HeatMap struct {
	Observer observer.Grid   // Observers providing a Grid for contours.
	Palette  palette.Palette // Color palette. Optional.
	Min      float64         // Minimum value for color mapping. Optional.
	Max      float64         // Maximum value for color mapping. Optional. Is set to 1.0 if both Min and Max are zero.
	Labels   Labels          // Labels for plot and axes. Optional.

	data  plotGrid
	scale float64
}

// Initialize the drawer.
func (s *HeatMap) Initialize(w *ecs.World, win *pixelgl.Window) {
	s.Observer.Initialize(w)
	s.scale = calcScaleCorrection()

	if s.Min == 0 && s.Max == 0 {
		s.Max = 1
	}
}

// Update the drawer.
func (s *HeatMap) Update(w *ecs.World) {
	s.Observer.Update(w)

	s.data = plotGrid{
		Grid: s.Observer,
	}
}

// UpdateInputs handles input events of the previous frame update.
func (s *HeatMap) UpdateInputs(w *ecs.World, win *pixelgl.Window) {}

// Draw the drawer.
func (s *HeatMap) Draw(w *ecs.World, win *pixelgl.Window) {
	s.updateData(w)

	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	if width <= 0 || height <= 0 {
		return
	}

	c := vgimg.New(vg.Points(width*s.scale)-10, vg.Points(height*s.scale)-10)

	p := plot.New()
	setLabels(p, s.Labels)

	p.X.Tick.Marker = removeLastTicks{}

	cols := s.Palette.Colors()
	heat := &plotter.HeatMap{
		GridXYZ:    &s.data,
		Palette:    s.Palette,
		Rasterized: false,
		Underflow:  cols[0],
		Overflow:   cols[len(cols)-1],
		Min:        s.Min,
		Max:        s.Max,
	}

	p.Add(heat)

	win.Clear(color.White)
	p.Draw(draw.New(c))

	img := c.Image()
	picture := pixel.PictureDataFromImage(img)

	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(picture.Rect.W()/2.0+5, picture.Rect.H()/2.0+5)))
}

func (s *HeatMap) updateData(w *ecs.World) {
	s.data.Values = s.Observer.Values(w)
}
