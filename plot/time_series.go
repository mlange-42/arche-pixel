package plot

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mlange-42/arche-model/observer"
	"github.com/mlange-42/arche-model/resource"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
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
//
// Expects a world resource of type [github.com/mlange-42/arche-model/resource.Tick].
type TimeSeries struct {
	Observer       observer.Row // Observer providing a data row per update.
	UpdateInterval int          // Interval for updating the observer, in model ticks. Optional.

	headers []string
	series  []plotter.XYs
	scale   float64
	tickRes generic.Resource[resource.Tick]
}

func (s *TimeSeries) append(x float64, values []float64) {
	for i := 0; i < len(s.series); i++ {
		s.series[i] = append(s.series[i], plotter.XY{X: x, Y: values[i]})
	}
}

// Initialize the drawer.
func (s *TimeSeries) Initialize(w *ecs.World, win *pixelgl.Window) {
	s.Observer.Initialize(w)

	s.tickRes = generic.NewResource[resource.Tick](w)

	s.headers = s.Observer.Header()
	s.series = make([]plotter.XYs, len(s.headers))

	s.scale = calcScaleCorrection()
}

// Update the drawer.
func (s *TimeSeries) Update(w *ecs.World) {
	tick := s.tickRes.Get().Tick

	if s.UpdateInterval <= 1 || tick%int64(s.UpdateInterval) == 0 {
		s.Observer.Update(w)
		s.append(float64(tick), s.Observer.Values(w))
	}
}

// Draw the drawer.
func (s *TimeSeries) Draw(w *ecs.World, win *pixelgl.Window) {
	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	c := vgimg.New(vg.Points(width*s.scale)-10, vg.Points(height*s.scale)-10)

	p := plot.New()
	p.X.Tick.Label.Font.Size = 12
	p.Y.Tick.Label.Font.Size = 12

	p.Legend = plot.NewLegend()

	for i := 0; i < len(s.series); i++ {
		lines, err := plotter.NewLine(s.series[i])
		if err != nil {
			panic(err)
		}
		lines.Color = defaultColors[i%len(defaultColors)]
		p.Add(lines)
		p.Legend.Add(s.headers[i], lines)
	}

	win.Clear(color.White)
	p.Draw(draw.New(c))

	img := c.Image()
	picture := pixel.PictureDataFromImage(img)

	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(picture.Rect.W()/2.0+5, picture.Rect.H()/2.0+5)))
}
