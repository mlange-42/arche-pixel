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

// Lines plot drawer.
//
// Creates a line series per column of the observer.
// Replaces the complete data by the table provided by the observer on every update.
// Particularly useful for live histograms.
type Lines struct {
	Observer observer.Table // Observer providing a data series as lines.
	X        string         // X column name. Optional. Defaults to row index.
	Y        []string       // Y column names. Optional. Defaults to all but X column.

	xIndex   int
	yIndices []int

	headers []string
	series  []plotter.XYs
	scale   float64
}

// Initialize the drawer.
func (l *Lines) Initialize(w *ecs.World, win *pixelgl.Window) {
	l.Observer.Initialize(w)

	l.headers = l.Observer.Header()

	l.scale = calcScaleCorrection()

	if l.X == "" {
		l.xIndex = -1
	} else {
		l.xIndex = -1
		for i, h := range l.headers {
			if h == l.X {
				l.xIndex = i
				break
			}
		}
		if l.xIndex < 0 {
			panic(fmt.Sprintf("x column '%s' not found", l.X))
		}
	}

	if len(l.Y) == 0 {
		l.yIndices = make([]int, 0, len(l.headers))
		for i := 0; i < len(l.headers); i++ {
			if i != l.xIndex {
				l.yIndices = append(l.yIndices, i)
			}
		}
	} else {
		l.yIndices = make([]int, len(l.Y))
		for i, y := range l.Y {
			idx := -1
			for j, h := range l.headers {
				if h == y {
					idx = j
					break
				}
			}
			if idx < 0 {
				panic(fmt.Sprintf("y column '%s' not found", y))
			}
			l.yIndices[i] = idx
		}
	}

	l.series = make([]plotter.XYs, len(l.yIndices))
}

// Update the drawer.
func (l *Lines) Update(w *ecs.World) {
	l.Observer.Update(w)
}

// UpdateInputs handles input events of the previous frame update.
func (l *Lines) UpdateInputs(w *ecs.World, win *pixelgl.Window) {}

// Draw the drawer.
func (l *Lines) Draw(w *ecs.World, win *pixelgl.Window) {
	l.updateData(w)

	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	c := vgimg.New(vg.Points(width*l.scale)-10, vg.Points(height*l.scale)-10)

	p := plot.New()
	p.X.Tick.Label.Font.Size = 12
	p.Y.Tick.Label.Font.Size = 12

	p.Legend = plot.NewLegend()

	for i := 0; i < len(l.series); i++ {
		idx := l.yIndices[i]
		lines, err := plotter.NewLine(l.series[i])
		if err != nil {
			panic(err)
		}
		lines.Color = defaultColors[i%len(defaultColors)]
		p.Add(lines)
		p.Legend.Add(l.headers[idx], lines)
	}

	win.Clear(color.White)
	p.Draw(draw.New(c))

	img := c.Image()
	picture := pixel.PictureDataFromImage(img)

	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(picture.Rect.W()/2.0+5, picture.Rect.H()/2.0+5)))
}

func (l *Lines) updateData(w *ecs.World) {
	data := l.Observer.Values(w)
	xi := l.xIndex
	yis := l.yIndices

	for i, idx := range yis {
		l.series[i] = l.series[i][:0]
		for j, row := range data {
			x := float64(j)
			if xi >= 0 {
				x = row[xi]
			}
			l.series[i] = append(l.series[i], plotter.XY{X: x, Y: row[idx]})
		}
	}
}
