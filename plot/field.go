package plot

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mlange-42/arche-model/observer"
	"github.com/mlange-42/arche/ecs"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

// Field plot drawer.
//
// Plots a vector field from a GridLayers observer.
type Field struct {
	Observer observer.GridLayers // Observers providing field component grids.
	Labels   Labels              // Labels for plot and axes. Optional.
	Layers   *[2]int             // Layer indices. Optional, defaults to (0, 1).

	data  plotField
	scale float64
}

// Initialize the drawer.
func (c *Field) Initialize(w *ecs.World, win *pixelgl.Window) {
	c.Observer.Initialize(w)

	c.data = plotField{
		GridLayers: c.Observer,
	}

	if c.Layers == nil {
		c.Layers = &[2]int{0, 1}
	}
	layers := c.Observer.Layers()
	for _, l := range c.Layers {
		if layers <= l {
			panic(fmt.Sprintf("layer index %d out of range", l))
		}
	}

	c.scale = calcScaleCorrection()
}

// Update the drawer.
func (c *Field) Update(w *ecs.World) {
	c.Observer.Update(w)
}

// UpdateInputs handles input events of the previous frame update.
func (c *Field) UpdateInputs(w *ecs.World, win *pixelgl.Window) {}

// Draw the drawer.
func (c *Field) Draw(w *ecs.World, win *pixelgl.Window) {
	c.updateData(w)

	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	if width <= 0 || height <= 0 {
		return
	}

	canvas := vgimg.New(vg.Points(width*c.scale)-10, vg.Points(height*c.scale)-10)

	p := plot.New()
	setLabels(p, c.Labels)

	p.X.Tick.Marker = removeLastTicks{}

	field := plotter.NewField(&c.data)

	p.Add(field)

	win.Clear(color.White)
	p.Draw(draw.New(canvas))

	img := canvas.Image()
	picture := pixel.PictureDataFromImage(img)

	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(picture.Rect.W()/2.0+5, picture.Rect.H()/2.0+5)))
}

func (c *Field) updateData(w *ecs.World) {
	values := c.Observer.Values(w)
	c.data.XValues = values[c.Layers[0]]
	c.data.YValues = values[c.Layers[1]]
}

type plotField struct {
	observer.GridLayers
	XValues []float64
	YValues []float64
}

func (f *plotField) Vector(c, r int) plotter.XY {
	w, _ := f.GridLayers.Dims()
	return plotter.XY{
		X: f.XValues[r*w+c],
		Y: f.YValues[r*w+c],
	}
}
