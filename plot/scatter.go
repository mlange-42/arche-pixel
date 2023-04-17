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
type Scatter struct {
	Observers []observer.Table // Observers providing XY data series.
	X         []string         // X column name per observer. Optional. Defaults to first column. Empty strings also falls back to the default.
	Y         []string         // Y column name per observer. Optional. Defaults to second column. Empty strings also falls back to the default.

	xIndices []int
	yIndices []int
	labels   []string

	series []plotter.XYs
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
	l.yIndices = make([]int, numObs)
	l.labels = make([]string, numObs)
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
		if len(l.Y) == 0 || l.Y[i] == "" {
			l.yIndices[i] = 1
		} else {
			l.yIndices[i], ok = find(header, l.Y[i])
			if !ok {
				panic(fmt.Sprintf("y column '%s' not found", l.Y[i]))
			}
		}
		l.labels[i] = header[l.yIndices[i]]
	}

	l.scale = calcScaleCorrection()
	l.series = make([]plotter.XYs, len(l.yIndices))
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

	p.Legend = plot.NewLegend()

	for i := 0; i < len(l.series); i++ {
		lines, err := plotter.NewScatter(l.series[i])
		if err != nil {
			panic(err)
		}
		lines.Color = defaultColors[i%len(defaultColors)]
		p.Add(lines)
		p.Legend.Add(l.labels[i], lines)
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
	yis := l.yIndices

	for i, yi := range yis {
		xi := xis[i]
		l.series[i] = l.series[i][:0]
		data := l.Observers[i].Values(w)
		for _, row := range data {
			l.series[i] = append(l.series[i], plotter.XY{X: row[xi], Y: row[yi]})
		}
	}
}
