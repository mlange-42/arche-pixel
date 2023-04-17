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

// Scatter plot drawer.
//
// Creates a scatter plot from multiple observers.
// Supports multiple series per observer. The series in a particular observer must share a common X column.
type Scatter struct {
	Observers []observer.Table // Observers providing XY data series.
	X         []string         // X column name per observer. Optional. Defaults to first column. Empty strings also falls back to the default.
	Y         [][]string       // Y column names per observer. Optional. Defaults to second column. Empty strings also falls back to the default.

	xIndices []int
	yIndices [][]int
	labels   [][]string

	series [][]plotter.XYs
	scale  float64
}

// Initialize the drawer.
func (s *Scatter) Initialize(w *ecs.World, win *pixelgl.Window) {
	numObs := len(s.Observers)
	if len(s.X) != 0 && len(s.X) != numObs {
		panic("length of X not equal to length of Observers")
	}
	if len(s.Y) != 0 && len(s.Y) != numObs {
		panic("length of Y not equal to length of Observers")
	}

	s.xIndices = make([]int, numObs)
	s.yIndices = make([][]int, numObs)
	s.labels = make([][]string, numObs)
	s.series = make([][]plotter.XYs, numObs)
	var ok bool
	for i := 0; i < numObs; i++ {
		obs := s.Observers[i]
		obs.Initialize(w)
		header := obs.Header()
		if len(s.X) == 0 || s.X[i] == "" {
			s.xIndices[i] = 0
		} else {
			s.xIndices[i], ok = find(header, s.X[i])
			if !ok {
				panic(fmt.Sprintf("x column '%s' not found", s.X[i]))
			}
		}
		if len(s.Y) == 0 || len(s.Y[i]) == 0 {
			s.yIndices[i] = []int{1}
			s.labels[i] = []string{header[1]}
			s.series[i] = make([]plotter.XYs, 1)
		} else {
			numY := len(s.Y[i])
			s.yIndices[i] = make([]int, numY)
			s.labels[i] = make([]string, numY)
			s.series[i] = make([]plotter.XYs, numY)
			for j, y := range s.Y[i] {
				idx, ok := find(header, y)
				if !ok {
					panic(fmt.Sprintf("y column '%s' not found", y))
				}
				s.yIndices[i][j] = idx
				s.labels[i][j] = header[idx]
			}

		}
	}

	s.scale = calcScaleCorrection()
}

// Update the drawer.
func (s *Scatter) Update(w *ecs.World) {
	for _, obs := range s.Observers {
		obs.Update(w)
	}
}

// UpdateInputs handles input events of the previous frame update.
func (s *Scatter) UpdateInputs(w *ecs.World, win *pixelgl.Window) {}

// Draw the drawer.
func (s *Scatter) Draw(w *ecs.World, win *pixelgl.Window) {
	s.updateData(w)

	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	if width <= 0 || height <= 0 {
		return
	}

	c := vgimg.New(vg.Points(width*s.scale)-10, vg.Points(height*s.scale)-10)

	p := plot.New()
	p.X.Tick.Label.Font.Size = 12
	p.Y.Tick.Label.Font.Size = 12
	p.Y.Tick.Label.Font.Variant = "Mono"
	p.X.Tick.Label.Font.Variant = "Mono"
	p.Y.Tick.Marker = paddedTicks{}

	p.Legend = plot.NewLegend()
	p.Legend.TextStyle.Font.Variant = "Mono"

	cnt := 0
	for i := 0; i < len(s.xIndices); i++ {
		ys := s.yIndices[i]
		for j := 0; j < len(ys); j++ {
			lines, err := plotter.NewScatter(s.series[i][j])
			if err != nil {
				panic(err)
			}
			lines.Color = defaultColors[cnt%len(defaultColors)]
			p.Add(lines)
			p.Legend.Add(s.labels[i][j], lines)
			cnt++
		}
	}

	win.Clear(color.White)
	p.Draw(draw.New(c))

	img := c.Image()
	picture := pixel.PictureDataFromImage(img)

	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(picture.Rect.W()/2.0+5, picture.Rect.H()/2.0+5)))
}

func (s *Scatter) updateData(w *ecs.World) {
	xis := s.xIndices

	for i := 0; i < len(xis); i++ {
		data := s.Observers[i].Values(w)
		xi := xis[i]
		ys := s.yIndices[i]
		for j := 0; j < len(ys); j++ {
			s.series[i][j] = s.series[i][j][:0]
			for _, row := range data {
				s.series[i][j] = append(s.series[i][j], plotter.XY{X: row[xi], Y: row[ys[j]]})
			}
		}
	}
}
