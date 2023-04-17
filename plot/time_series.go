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

// TimeSeries plot drawer.
//
// Creates a line series per column of the observer.
// Adds one row to the data per update.
type TimeSeries struct {
	Observer       observer.Row // Observer providing a data row per update.
	UpdateInterval int          // Interval for getting data from the the observer, in model ticks. Optional.
	Labels         Labels       // Labels for plot and axes. Optional.

	headers []string
	series  []plotter.XYs
	scale   float64
	step    int64
}

// append a y value to each series, with a common x value.
func (t *TimeSeries) append(x float64, values []float64) {
	for i := 0; i < len(t.series); i++ {
		t.series[i] = append(t.series[i], plotter.XY{X: x, Y: values[i]})
	}
}

// Initialize the drawer.
func (t *TimeSeries) Initialize(w *ecs.World, win *pixelgl.Window) {
	t.Observer.Initialize(w)

	t.headers = t.Observer.Header()
	t.series = make([]plotter.XYs, len(t.headers))

	t.scale = calcScaleCorrection()
	t.step = 0
}

// Update the drawer.
func (t *TimeSeries) Update(w *ecs.World) {
	t.Observer.Update(w)
	if t.UpdateInterval <= 1 || t.step%int64(t.UpdateInterval) == 0 {
		t.append(float64(t.step), t.Observer.Values(w))
	}
	t.step++
}

// UpdateInputs handles input events of the previous frame update.
func (t *TimeSeries) UpdateInputs(w *ecs.World, win *pixelgl.Window) {}

// Draw the drawer.
func (t *TimeSeries) Draw(w *ecs.World, win *pixelgl.Window) {
	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	if width <= 0 || height <= 0 {
		return
	}

	c := vgimg.New(vg.Points(width*t.scale)-10, vg.Points(height*t.scale)-10)

	p := plot.New()
	setLabels(p, t.Labels)

	p.X.Tick.Marker = removeLastTicks{}

	p.Legend = plot.NewLegend()
	p.Legend.TextStyle.Font.Variant = "Mono"

	for i := 0; i < len(t.series); i++ {
		lines, err := plotter.NewLine(t.series[i])
		if err != nil {
			panic(err)
		}
		lines.Color = defaultColors[i%len(defaultColors)]
		p.Add(lines)
		p.Legend.Add(t.headers[i], lines)
	}

	win.Clear(color.White)
	p.Draw(draw.New(c))

	img := c.Image()
	picture := pixel.PictureDataFromImage(img)

	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(picture.Rect.W()/2.0+5, picture.Rect.H()/2.0+5)))
}
