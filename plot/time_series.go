package plot

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mlange-42/arche-model/observer"
	"github.com/mlange-42/arche/ecs"
	"golang.org/x/image/colornames"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

var defaultColors = []color.Color{
	colornames.Blue,
	colornames.Orange,
	colornames.Green,
	colornames.Purple,
	colornames.Red,
	colornames.Turquoise,
}

// TimeSeries plot drawer.
type TimeSeries struct {
	Observer       observer.Row // Observer providing a data row per update.
	UpdateInterval int          // Interval for updating the observer, in model ticks. Optional.

	headers []string
	series  []plotter.XYs
	scale   float64
	step    int64
}

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
	if t.UpdateInterval <= 1 || t.step%int64(t.UpdateInterval) == 0 {
		t.Observer.Update(w)
		t.append(float64(t.step), t.Observer.Values(w))
	}
	t.step++
}

// Draw the drawer.
func (t *TimeSeries) Draw(w *ecs.World, win *pixelgl.Window) {
	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	c := vgimg.New(vg.Points(width*t.scale)-10, vg.Points(height*t.scale)-10)

	p := plot.New()
	p.X.Tick.Label.Font.Size = 12
	p.Y.Tick.Label.Font.Size = 12

	p.Legend = plot.NewLegend()

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
