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
func (l *Scatter) Initialize(w *ecs.World, win *pixelgl.Window) {
	numObs := len(l.Observers)
	if len(l.X) != 0 && len(l.X) != numObs {
		panic("length of X not equal to length of Observers")
	}
	if len(l.Y) != 0 && len(l.Y) != numObs {
		panic("length of Y not equal to length of Observers")
	}

	l.xIndices = make([]int, numObs)
	l.yIndices = make([][]int, numObs)
	l.labels = make([][]string, numObs)
	l.series = make([][]plotter.XYs, numObs)
	var ok bool
	for i := 0; i < numObs; i++ {
		obs := l.Observers[i]
		obs.Initialize(w)
		header := obs.Header()
		if len(l.X) == 0 || l.X[i] == "" {
			l.xIndices[i] = 0
		} else {
			l.xIndices[i], ok = find(header, l.X[i])
			if !ok {
				panic(fmt.Sprintf("x column '%s' not found", l.X[i]))
			}
		}
		if len(l.Y) == 0 || len(l.Y[i]) == 0 {
			l.yIndices[i] = []int{1}
			l.labels[i] = []string{header[1]}
			l.series[i] = make([]plotter.XYs, 1)
		} else {
			numY := len(l.Y[i])
			l.yIndices[i] = make([]int, numY)
			l.labels[i] = make([]string, numY)
			l.series[i] = make([]plotter.XYs, numY)
			for j, y := range l.Y[i] {
				idx, ok := find(header, y)
				if !ok {
					panic(fmt.Sprintf("y column '%s' not found", y))
				}
				l.yIndices[i][j] = idx
				l.labels[i][j] = header[idx]
			}

		}
	}

	l.scale = calcScaleCorrection()
}

// Update the drawer.
func (l *Scatter) Update(w *ecs.World) {
	for _, obs := range l.Observers {
		obs.Update(w)
	}
}

// UpdateInputs handles input events of the previous frame update.
func (l *Scatter) UpdateInputs(w *ecs.World, win *pixelgl.Window) {}

// Draw the drawer.
func (l *Scatter) Draw(w *ecs.World, win *pixelgl.Window) {
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

	p.Legend = plot.NewLegend()
	p.Legend.TextStyle.Font.Variant = "Mono"

	cnt := 0
	for i := 0; i < len(l.xIndices); i++ {
		ys := l.yIndices[i]
		for j := 0; j < len(ys); j++ {
			lines, err := plotter.NewScatter(l.series[i][j])
			if err != nil {
				panic(err)
			}
			lines.Color = defaultColors[cnt%len(defaultColors)]
			p.Add(lines)
			p.Legend.Add(l.labels[i][j], lines)
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

func (l *Scatter) updateData(w *ecs.World) {
	xis := l.xIndices

	for i := 0; i < len(xis); i++ {
		data := l.Observers[i].Values(w)
		xi := xis[i]
		ys := l.yIndices[i]
		for j := 0; j < len(ys); j++ {
			l.series[i][j] = l.series[i][j][:0]
			for _, row := range data {
				l.series[i][j] = append(l.series[i][j], plotter.XY{X: row[xi], Y: row[ys[j]]})
			}
		}
	}
}
