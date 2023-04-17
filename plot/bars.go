package plot

import (
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

// Bars plot drawer.
//
// Creates a bar per column of the observer.
type Bars struct {
	Observer observer.Row // Observer providing a data series for bars.

	headers []string
	series  plotter.Values
	scale   float64
}

// Initialize the drawer.
func (l *Bars) Initialize(w *ecs.World, win *pixelgl.Window) {
	l.Observer.Initialize(w)

	l.headers = l.Observer.Header()
	l.series = make([]float64, len(l.headers))

	l.scale = calcScaleCorrection()
}

// Update the drawer.
func (l *Bars) Update(w *ecs.World) {
	l.Observer.Update(w)
}

// UpdateInputs handles input events of the previous frame update.
func (l *Bars) UpdateInputs(w *ecs.World, win *pixelgl.Window) {}

// Draw the drawer.
func (l *Bars) Draw(w *ecs.World, win *pixelgl.Window) {
	l.updateData(w)

	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	c := vgimg.New(vg.Points(width*l.scale)-10, vg.Points(height*l.scale)-10)

	p := plot.New()
	p.X.Tick.Label.Font.Size = 12
	p.Y.Tick.Label.Font.Size = 12
	p.Y.Tick.Label.Font.Variant = "Mono"
	p.X.Tick.Label.Font.Variant = "Mono"
	p.Y.Tick.Marker = paddedTicks{}

	bw := 0.5 * (width - 50) / float64(len(l.series))
	bars, err := plotter.NewBarChart(l.series, vg.Points(bw))
	if err != nil {
		panic(err)
	}
	bars.Color = defaultColors[0]
	p.Add(bars)
	p.NominalX(l.headers...)

	win.Clear(color.White)
	p.Draw(draw.New(c))

	img := c.Image()
	picture := pixel.PictureDataFromImage(img)

	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(picture.Rect.W()/2.0+5, picture.Rect.H()/2.0+5)))
}

func (l *Bars) updateData(w *ecs.World) {
	l.series = l.Observer.Values(w)
}
