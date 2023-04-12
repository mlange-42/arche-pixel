package plot

import (
	"fmt"
	"image/color"
	"time"

	px "github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/mlange-42/arche/ecs"
)

// PerfStats drawer for performance statistics.
//
// Adds an overlay with performance statistics in the top left corner of the window.
type PerfStats struct {
	SampleInterval time.Duration // Approx. time between measurements. Optional, default 1 second.
	drawer         imdraw.IMDraw
	stats          tempStats
	summary        *text.Text
	frameTimer     frameTimer
	startTime      time.Time
	lastPlotUpdate time.Time
	step           int64
}

// Initialize the system
func (m *PerfStats) Initialize(w *ecs.World, win *pixelgl.Window) {
	if m.SampleInterval <= 0 {
		m.SampleInterval = time.Second
	}
	m.lastPlotUpdate = time.Now()
	m.startTime = m.lastPlotUpdate

	m.drawer = *imdraw.New(nil)

	m.summary = text.New(px.V(0, 0), font)

	m.step = 0

	st := w.Stats()
	m.stats.Entities = st.Entities.Used
	m.stats.Mem = st.Memory
}

// Update the drawer.
func (m *PerfStats) Update(w *ecs.World) {
	t := time.Now()
	m.frameTimer.Update(m.step, t)

	if t.Sub(m.lastPlotUpdate) >= m.SampleInterval {
		st := w.Stats()
		m.stats.Entities = st.Entities.Used
		m.stats.Mem = st.Memory
		m.lastPlotUpdate = t
	}

	m.step++
}

// UpdateInputs handles input events of the previous frame update.
func (m *PerfStats) UpdateInputs(w *ecs.World, win *pixelgl.Window) {}

// Draw the system
func (m *PerfStats) Draw(w *ecs.World, win *pixelgl.Window) {
	m.summary.Clear()
	mem, units := toMemText(m.stats.Mem)
	fmt.Fprintf(
		m.summary, "Tick: %7d\nEnt.: %7d\nTPS: %8.1f\nTPT: %6.2fms\nMem: %6.1f%s\nTime: %7s",
		m.step, m.stats.Entities, m.frameTimer.FPS(),
		float64(m.frameTimer.FrameTime().Microseconds())/1000,
		mem, units, time.Since(m.startTime).Round(time.Second),
	)

	dr := &m.drawer
	height := win.Canvas().Bounds().H()
	x0 := 10.0
	y0 := height - 20.0

	v1 := px.V(x0+m.summary.Bounds().Min.X-5, y0+m.summary.Bounds().Min.Y-5)
	v2 := px.V(x0+m.summary.Bounds().Max.X+5, y0+m.summary.Bounds().Max.Y+5)

	dr.Color = color.Black
	dr.Push(v1, v2)
	dr.Rectangle(0)

	dr.Color = color.White
	dr.Push(v1, v2)
	dr.Rectangle(1)

	dr.Draw(win)
	dr.Reset()
	dr.Clear()

	m.summary.Draw(win, px.IM.Moved(px.V(x0, y0)))

}

type tempStats struct {
	Entities int
	Mem      int
}
