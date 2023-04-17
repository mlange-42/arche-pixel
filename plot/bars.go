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
func (b *Bars) Initialize(w *ecs.World, win *pixelgl.Window) {
	b.Observer.Initialize(w)

	b.headers = b.Observer.Header()
	b.series = make([]float64, len(b.headers))

	b.scale = calcScaleCorrection()
}

// Update the drawer.
func (b *Bars) Update(w *ecs.World) {
	b.Observer.Update(w)
}

// UpdateInputs handles input events of the previous frame update.
func (b *Bars) UpdateInputs(w *ecs.World, win *pixelgl.Window) {}

// Draw the drawer.
func (b *Bars) Draw(w *ecs.World, win *pixelgl.Window) {
	b.updateData(w)

	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	if width <= 0 || height <= 0 {
		return
	}

	c := vgimg.New(vg.Points(width*b.scale)-10, vg.Points(height*b.scale)-10)

	p := plot.New()
	p.X.Tick.Label.Font.Size = 12
	p.Y.Tick.Label.Font.Size = 12
	p.Y.Tick.Label.Font.Variant = "Mono"
	p.X.Tick.Label.Font.Variant = "Mono"
	p.Y.Tick.Marker = paddedTicks{}

	bw := 0.5 * (width - 50) / float64(len(b.series))
	bars, err := plotter.NewBarChart(b.series, vg.Points(bw))
	if err != nil {
		panic(err)
	}
	bars.Color = defaultColors[0]
	p.Add(bars)
	p.NominalX(b.headers...)

	win.Clear(color.White)
	p.Draw(draw.New(c))

	img := c.Image()
	picture := pixel.PictureDataFromImage(img)

	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(picture.Rect.W()/2.0+5, picture.Rect.H()/2.0+5)))
}

func (b *Bars) updateData(w *ecs.World) {
	b.series = b.Observer.Values(w)
}
