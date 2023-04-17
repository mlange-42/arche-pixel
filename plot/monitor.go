package plot

import (
	"fmt"
	"image/color"
	"math"
	"strings"
	"time"

	px "github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/mlange-42/arche-pixel/window"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/stats"
)

type timeSeriesType uint8

const (
	tsEntities timeSeriesType = iota
	tsEntityCap
	tsMemory
	tsTickPerSec
	tsLast
)

// NewMonitorWindow creates a window with [Monitor] drawer, for immediate use as a system.
//
// Also adds [Controls] for pausing/resuming the simulation.
func NewMonitorWindow(drawInterval int) *window.Window {
	return (&window.Window{
		Title:        "Monitor",
		DrawInterval: drawInterval,
	}).With(
		&Monitor{
			SampleInterval: time.Second / 3,
		},
		&Controls{},
	)
}

// Monitor drawer for visualizing world and performance statistics.
type Monitor struct {
	PlotCapacity   int           // Number of values in time series plots. Optional, default 300.
	SampleInterval time.Duration // Approx. time between measurements for time series plots. Optional, default 1 second.
	HidePlots      bool          // Hides time series plots
	HideArchetypes bool          // Hides archetype stats
	scale          float64
	drawer         imdraw.IMDraw
	summary        *text.Text
	timeSeries     timeSeries
	frameTimer     frameTimer
	archetypes     archetypes
	text           *text.Text
	startTime      time.Time
	lastPlotUpdate time.Time
	step           int64
}

// Initialize the system
func (m *Monitor) Initialize(w *ecs.World, win *pixelgl.Window) {
	if m.PlotCapacity <= 0 {
		m.PlotCapacity = 300
	}
	if m.SampleInterval <= 0 {
		m.SampleInterval = time.Second
	}
	m.lastPlotUpdate = time.Now()
	m.startTime = m.lastPlotUpdate

	m.drawer = *imdraw.New(nil)

	m.scale = calcScaleCorrection()

	m.summary = text.New(px.V(0, 0), defaultFont)

	m.timeSeries = newTimeSeries(m.PlotCapacity)
	for i := 0; i < len(m.timeSeries.Text); i++ {
		m.timeSeries.Text[i] = text.New(px.V(0, 0), defaultFont)
	}
	fmt.Fprintf(m.timeSeries.Text[tsEntities], "Entities")
	fmt.Fprintf(m.timeSeries.Text[tsEntityCap], "Capacity")
	fmt.Fprintf(m.timeSeries.Text[tsMemory], "Memory")
	fmt.Fprintf(m.timeSeries.Text[tsTickPerSec], "TPS")

	m.text = text.New(px.V(0, 0), defaultFont)
	m.text.Color = color.RGBA{200, 200, 200, 255}

	m.step = 0
}

// Update the drawer.
func (m *Monitor) Update(w *ecs.World) {
	t := time.Now()
	m.frameTimer.Update(m.step, t)

	if !m.HidePlots && t.Sub(m.lastPlotUpdate) >= m.SampleInterval {
		stats := w.Stats()
		m.archetypes.Update(stats)
		m.timeSeries.append(
			stats.Entities.Used, stats.Entities.Total, stats.Memory,
			int(m.frameTimer.FPS()*1000000),
		)
		m.lastPlotUpdate = t
	}
	m.step++
}

// UpdateInputs handles input events of the previous frame update.
func (m *Monitor) UpdateInputs(w *ecs.World, win *pixelgl.Window) {}

// Draw the system
func (m *Monitor) Draw(w *ecs.World, win *pixelgl.Window) {
	stats := w.Stats()
	m.archetypes.Update(stats)

	m.summary.Clear()
	mem, units := toMemText(stats.Memory)
	fmt.Fprintf(
		m.summary, "Tick: %7d  |  Entities: %7d  |  Comp: %3d  |  Arch: %3d  |  Mem: %6.1f %s  |  TPS: %6.1f | TPT: %6.2f ms | Time: %s",
		m.step, stats.Entities.Used, stats.ComponentCount, len(stats.Archetypes), mem, units, m.frameTimer.FPS(),
		float64(m.frameTimer.FrameTime().Microseconds())/1000, time.Since(m.startTime).Round(time.Second),
	)

	numArch := len(stats.Archetypes)
	maxCapacity := 0
	for i := 0; i < numArch; i++ {
		cap := stats.Archetypes[i].Capacity
		if cap > maxCapacity {
			maxCapacity = cap
		}
	}
	dr := &m.drawer
	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()
	x0 := 10.0
	y0 := height - 20.0

	m.summary.Draw(win, px.IM.Moved(px.V(x0, y0)))
	y0 -= 10

	if !m.HidePlots {
		plotY0 := y0
		plotHeight := (height - 60) / 3
		if plotHeight > 150 {
			plotHeight = 150
		}
		plotWidth := (width - 20) * 0.25
		if m.HideArchetypes {
			plotWidth = width - 20
		}
		m.drawPlot(win, x0, plotY0-plotHeight, plotWidth, plotHeight, tsEntities, tsEntityCap)
		plotY0 -= plotHeight + 10
		m.drawPlot(win, x0, plotY0-plotHeight, plotWidth, plotHeight, tsMemory)
		plotY0 -= plotHeight + 10
		m.drawPlot(win, x0, plotY0-plotHeight, plotWidth, plotHeight, tsTickPerSec)

		x0 += math.Ceil(plotWidth + 10)
	}

	archHeight := math.Ceil((y0 - 10) / float64(numArch+1))
	if !m.HideArchetypes {
		if archHeight >= 8 {
			archWidth := width - x0 - 10
			if archHeight > 20 {
				archHeight = 20
			}
			m.drawArchetypeScales(
				win, x0, y0-archHeight, archWidth, archHeight, maxCapacity,
			)
			for i := 0; i < numArch; i++ {
				m.drawArchetype(
					win, x0, y0-float64(i+2)*archHeight, archWidth, archHeight,
					maxCapacity, &stats.Archetypes[i], m.archetypes.Components[i],
				)
			}
		} else {
			m.text.Clear()
			fmt.Fprintf(m.text, "Too many archetypes")
			m.text.Draw(win, px.IM.Moved(px.V(x0, y0-10)))
		}
	}

	dr.Draw(win)
	dr.Clear()
}

func (m *Monitor) drawArchetypeScales(win *pixelgl.Window, x, y, w, h float64, max int) {
	dr := &m.drawer
	step := calcTicksStep(float64(max), 8)
	if step < 1 {
		return
	}
	drawStep := w * step / float64(max)

	dr.Color = color.RGBA{140, 140, 140, 255}
	dr.Push(px.V(x, y+2), px.V(x+w, y+2))
	dr.Line(1)
	dr.Reset()

	steps := int(float64(max) / step)
	for i := 0; i <= steps; i++ {
		xi := float64(i)
		dr.Push(px.V(x+xi*drawStep, y+2), px.V(x+xi*drawStep, y+7))
		dr.Line(1)
		dr.Reset()

		val := i * int(step)
		m.text.Clear()
		fmt.Fprintf(m.text, "%d", val)
		m.text.Draw(win, px.IM.Moved(px.V(math.Floor(x+xi*drawStep-m.text.Bounds().W()/2), y+10)))
	}
}

func (m *Monitor) drawArchetype(win *pixelgl.Window, x, y, w, h float64, max int, arch *stats.ArchetypeStats, text *text.Text) {
	dr := &m.drawer

	cap := float64(arch.Capacity) / float64(max)
	cnt := float64(arch.Size) / float64(max)

	//dr.Color = color.RGBA{160, 40, 40, 255}
	dr.Color = color.RGBA{0, 120, 40, 255}
	dr.Push(px.V(x, y), px.V(x+w*cnt, y+h))
	dr.Rectangle(0)
	dr.Reset()

	//dr.Color = color.RGBA{40, 40, 160, 255}
	dr.Color = color.RGBA{20, 60, 25, 255}
	dr.Push(px.V(x+w*cnt, y), px.V(x+w*cap, y+h))
	dr.Rectangle(0)
	dr.Reset()

	dr.Color = color.RGBA{40, 40, 40, 255}
	dr.Push(px.V(x, y), px.V(x+w, y+h))
	dr.Rectangle(1)
	dr.Reset()

	dr.Draw(win)
	dr.Clear()

	text.Draw(win, px.IM.Moved(px.V(x+3, y+3)))
}

func (m *Monitor) drawPlot(win *pixelgl.Window, x, y, w, h float64, series ...timeSeriesType) {
	dr := &m.drawer

	dr.Color = color.RGBA{0, 25, 10, 255}
	dr.Push(px.V(x, y), px.V(x+w, y+h))
	dr.Rectangle(0)
	dr.Reset()

	yMax := 0
	for _, series := range series {
		values := m.timeSeries.Values[series]
		l := values.Len()
		for i := 0; i < l; i++ {
			v := values.Get(i)
			if v > yMax {
				yMax = v
			}
		}
	}

	dr.Color = color.White
	for _, series := range series {
		values := m.timeSeries.Values[series]
		numValues := values.Len()
		if numValues > 0 {
			xStep := w / float64(numValues-1)
			yScale := 0.95 * h / float64(yMax)

			for i := 0; i < numValues-1; i++ {
				xi := float64(i)
				x1, x2 := xi*xStep, xi*xStep+xStep
				y1, y2 := float64(values.Get(i))*yScale, float64(values.Get(i+1))*yScale

				dr.Push(px.V(x+x1, y+y1), px.V(x+x2, y+y2))
				dr.Line(1)
				dr.Reset()
			}
		}
	}

	dr.Color = color.RGBA{140, 140, 140, 255}
	dr.Push(px.V(x, y), px.V(x+w, y+h))
	dr.Rectangle(1)
	dr.Reset()

	dr.Draw(win)
	dr.Clear()

	if len(series) > 0 {
		text := m.timeSeries.Text[series[0]]
		text.Draw(win, px.IM.Moved(px.V(x+w-text.Bounds().W()-3, y+3)))
	}
}

func toMemText(bytes int) (float64, string) {
	if bytes < 10*1_048_576 {
		return float64(bytes) / 1024, "kB"
	}
	return float64(bytes) / 1_048_576, "MB"
}

type timeSeries struct {
	Values [tsLast]ringBuffer[int]
	Text   [tsLast]*text.Text
}

func newTimeSeries(cap int) timeSeries {
	ts := timeSeries{}
	for i := 0; i < int(tsLast); i++ {
		ts.Values[i] = newRingBuffer[int](cap)
	}
	return ts
}

func (t *timeSeries) append(entities, entityCap, memory, tps int) {
	t.Values[tsEntities].Add(entities)
	t.Values[tsEntityCap].Add(entityCap)
	t.Values[tsMemory].Add(memory)
	t.Values[tsTickPerSec].Add(tps)
}

type frameTimer struct {
	lastTick  int64
	lastTime  time.Time
	frameTime time.Duration
}

func (t *frameTimer) Update(tick int64, tm time.Time) {
	delta := tm.Sub(t.lastTime)

	if delta < time.Second {
		return
	}

	ticks := tick - t.lastTick

	if ticks > 0 {
		t.frameTime = delta / time.Duration(ticks)
	}

	t.lastTick = tick
	t.lastTime = tm
}

func (t *frameTimer) FrameTime() time.Duration {
	return t.frameTime
}

func (t *frameTimer) FPS() float64 {
	if t.frameTime == 0 {
		return 0
	}
	return 1.0 / t.frameTime.Seconds()
}

type archetypes struct {
	Components []*text.Text
}

func (a *archetypes) Update(stats *stats.WorldStats) {
	oldLen := len(a.Components)
	newLen := len(stats.Archetypes)

	if newLen == oldLen {
		return
	}

	for i := oldLen; i < newLen; i++ {
		text := text.New(px.V(0, 0), defaultFont)
		text.Color = color.RGBA{200, 200, 200, 255}
		arch := &stats.Archetypes[i]
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("%4d B: ", arch.MemoryPerEntity))
		types := arch.ComponentTypes
		for j := 0; j < len(types); j++ {
			sb.WriteString(types[j].Name())
			sb.WriteString(" ")
		}
		text.WriteString(sb.String())
		a.Components = append(a.Components, text)
	}
}
