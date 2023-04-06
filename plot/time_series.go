package plot

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mlange-42/arche-model/observer"
	"github.com/mlange-42/arche-pixel/window"
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

// TimeSeries plot reporter.
//
// If the world contains a resource of type [github.com/mlange-42/arche-model/resource.Termination],
// the model is terminated when the window is closed.
type TimeSeries struct {
	Bounds         window.Bounds
	Observer       observer.Row
	UpdateInterval int
	DrawInterval   int
	window         window.Window
	drawer         timeSeriesDrawer
	step           int64
}

// Initialize the system.
func (s *TimeSeries) Initialize(w *ecs.World) {
	s.step = 0
}

// InitializeUI the system.
func (s *TimeSeries) InitializeUI(w *ecs.World) {
	s.Observer.Initialize(w)

	s.drawer = timeSeriesDrawer{}
	s.drawer.addSeries(s.Observer.Header())

	s.window.DrawInterval = s.DrawInterval
	s.window.Bounds = s.Bounds
	s.window.Drawers = append([]window.Drawer{&s.drawer}, s.window.Drawers...)
	s.window.InitializeUI(w)
}

// Update the system.
func (s *TimeSeries) Update(w *ecs.World) {
	if s.UpdateInterval <= 1 || s.step%int64(s.UpdateInterval) == 0 {
		s.Observer.Update(w)
		s.drawer.append(float64(s.step), s.Observer.Values(w))
	}
	s.step++
}

// UpdateUI the system.
func (s *TimeSeries) UpdateUI(w *ecs.World) {
	s.window.UpdateUI(w)
}

// PostUpdateUI updates the GL window.
func (s *TimeSeries) PostUpdateUI(w *ecs.World) {
	s.window.PostUpdateUI(w)
}

// Finalize the system.
func (s *TimeSeries) Finalize(w *ecs.World) {}

// FinalizeUI the system.
func (s *TimeSeries) FinalizeUI(w *ecs.World) {
	s.window.FinalizeUI(w)
}

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
