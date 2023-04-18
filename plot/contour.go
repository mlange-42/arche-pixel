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

// Contour plot drawer.
//
// Plots a grid as a contours.
// For large grids, this is relatively slow.
// Consider using [Image] instead.ü
type Contour struct {
	Observer observer.Grid   // Observers providing a Grid for contours.
	Levels   []float64       // Levels for iso lines. Optional.
	Palette  palette.Palette // Color palette. Optional.
	Labels   Labels          // Labels for plot and axes. Optional.

	data  plotGrid
	scale float64
}

// Initialize the drawer.
func (s *Contour) Initialize(w *ecs.World, win *pixelgl.Window) {
	s.Observer.Initialize(w)
	s.scale = calcScaleCorrection()
}

// Update the drawer.
func (s *Contour) Update(w *ecs.World) {
	s.Observer.Update(w)

	s.data = plotGrid{
		Grid: s.Observer,
	}
}

// UpdateInputs handles input events of the previous frame update.
func (s *Contour) UpdateInputs(w *ecs.World, win *pixelgl.Window) {}

// Draw the drawer.
func (s *Contour) Draw(w *ecs.World, win *pixelgl.Window) {
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

	//contours := plotter.NewContour(&s.data, s.Levels, s.Palette)
	cols := s.Palette.Colors()
	min := 0.0
	max := 1.0
	if len(s.Levels) > 0 {
		min = s.Levels[0]
		max = s.Levels[len(s.Levels)-1]
	} else {
		s.Levels = []float64{0.01, 0.05, 0.25, 0.5, 0.75, 0.95, 0.99}
	}

	contours := &plotter.Contour{
		GridXYZ:    &s.data,
		Levels:     s.Levels,
		LineStyles: []draw.LineStyle{plotter.DefaultLineStyle},
		Palette:    s.Palette,
		Underflow:  cols[0],
		Overflow:   cols[len(cols)-1],
		Min:        min,
		Max:        max,
	}

	p.Add(contours)

	win.Clear(color.White)
	p.Draw(draw.New(c))

	img := c.Image()
	picture := pixel.PictureDataFromImage(img)

	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(picture.Rect.W()/2.0+5, picture.Rect.H()/2.0+5)))
}

func (s *Contour) updateData(w *ecs.World) {
	s.data.Values = s.Observer.Values(w)
}
