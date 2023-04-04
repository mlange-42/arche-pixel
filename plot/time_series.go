package plot

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-pixel/window"
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

// TimeSeries plot reporter.
type TimeSeries struct {
	Bounds         window.Bounds
	Observer       model.Observer
	UpdateInterval int
	DrawInterval   int
	window.Window
	drawer  timeSeriesDrawer
	timeRes generic.Resource[model.Time]
}

// Initialize the system.
func (s *TimeSeries) Initialize(w *ecs.World) {
	s.timeRes = generic.NewResource[model.Time](w)
}

// InitializeUI the system.
func (s *TimeSeries) InitializeUI(w *ecs.World) {
	s.Observer.Initialize(w)

	s.drawer = timeSeriesDrawer{}
	s.drawer.addSeries(s.Observer.Header(w))

	s.Window.DrawInterval = s.DrawInterval
	s.Window.Bounds = s.Bounds
	s.Window.Add(&s.drawer)
	s.Window.InitializeUI(w)
}

// Update the system.
func (s *TimeSeries) Update(w *ecs.World) {
	time := s.timeRes.Get()

	if s.UpdateInterval > 1 && time.Tick%int64(s.UpdateInterval) != 0 {
		return
	}
	s.Observer.Update(w)
	s.drawer.append(float64(time.Tick), s.Observer.Values(w))
}

// Finalize the system.
func (s *TimeSeries) Finalize(w *ecs.World) {}

type timeSeriesDrawer struct {
	headers []string
	series  []plotter.XYs
	scale   float64
}

func (s *timeSeriesDrawer) addSeries(names []string) {
	s.headers = names
	s.series = make([]plotter.XYs, len(s.headers))
}

func (s *timeSeriesDrawer) append(x float64, values []float64) {
	for i := 0; i < len(s.series); i++ {
		s.series[i] = append(s.series[i], plotter.XY{X: x, Y: values[i]})
	}
}

// Initialize the system
func (s *timeSeriesDrawer) Initialize(w *ecs.World, win *pixelgl.Window) {
	width := 100.0
	c := vgimg.New(vg.Points(width), vg.Points(width))
	img := c.Image()
	s.scale = width / float64(img.Bounds().Dx())
}

// Draw the system
func (s *timeSeriesDrawer) Draw(w *ecs.World, win *pixelgl.Window) {
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
