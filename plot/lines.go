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
type Lines struct {
	Observer       observer.Table // Observer providing a data series as lines.
	X              string         // X column name. Optional. Default first column.
	Y              []string       // Y column names. Optional. Default all but X column.
	UpdateInterval int            // Interval for updating the observer, in model ticks. Optional.

	xIndex   int
	yIndices []int

	headers []string
	series  []plotter.XYs
	scale   float64
	step    int64
}

// Initialize the drawer.
func (t *Lines) Initialize(w *ecs.World, win *pixelgl.Window) {
	t.Observer.Initialize(w)

	t.headers = t.Observer.Header()

	t.scale = calcScaleCorrection()
	t.step = 0

	if t.X == "" {
		t.xIndex = 0
	} else {
		t.xIndex = -1
		for i, h := range t.headers {
			if h == t.X {
				t.xIndex = i
				break
			}
		}
		if t.xIndex < 0 {
			panic(fmt.Sprintf("x column '%s' not found", t.X))
		}
	}

	if len(t.Y) == 0 {
		t.yIndices = make([]int, len(t.headers)-1)
		idx := 0
		for i := 0; i < len(t.headers); i++ {
			if i != t.xIndex {
				t.yIndices[idx] = i
				idx++
			}
		}
	} else {
		t.yIndices = make([]int, len(t.Y))
		for i, y := range t.Y {
			idx := -1
			for j, h := range t.headers {
				if h == y {
					idx = j
					break
				}
			}
			if idx < 0 {
				panic(fmt.Sprintf("y column '%s' not found", y))
			}
			t.yIndices[i] = idx
		}
	}

	t.series = make([]plotter.XYs, len(t.yIndices))
}

// Update the drawer.
func (t *Lines) Update(w *ecs.World) {
	t.Observer.Update(w)
	if t.UpdateInterval <= 1 || t.step%int64(t.UpdateInterval) == 0 {
		data := t.Observer.Values(w)
		x := t.xIndex
		ys := t.yIndices

		for i, idx := range ys {
			t.series[i] = t.series[i][:0]
			for _, row := range data {
				t.series[i] = append(t.series[i], plotter.XY{X: row[x], Y: row[idx]})
			}
		}
	}
	t.step++
}

// UpdateInputs handles input events of the previous frame update.
func (t *Lines) UpdateInputs(w *ecs.World, win *pixelgl.Window) {}

// Draw the drawer.
func (t *Lines) Draw(w *ecs.World, win *pixelgl.Window) {
	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	c := vgimg.New(vg.Points(width*t.scale)-10, vg.Points(height*t.scale)-10)

	p := plot.New()
	p.X.Tick.Label.Font.Size = 12
	p.Y.Tick.Label.Font.Size = 12

	p.Legend = plot.NewLegend()

	for i := 0; i < len(t.series); i++ {
		idx := t.yIndices[i]
		lines, err := plotter.NewLine(t.series[i])
		if err != nil {
			panic(err)
		}
		lines.Color = defaultColors[i%len(defaultColors)]
		p.Add(lines)
		p.Legend.Add(t.headers[idx], lines)
	}

	win.Clear(color.White)
	p.Draw(draw.New(c))

	img := c.Image()
	picture := pixel.PictureDataFromImage(img)

	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(win, pixel.IM.Moved(pixel.V(picture.Rect.W()/2.0+5, picture.Rect.H()/2.0+5)))
}
