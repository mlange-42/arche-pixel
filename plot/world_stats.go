package plot

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	px "github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/mlange-42/arche-model/resource"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/ecs/stats"
	"github.com/mlange-42/arche/generic"
	"golang.org/x/image/font/basicfont"
)

type timeSeriesType uint8

const (
	tsEntities timeSeriesType = iota
	tsEntityCap
	tsMemory
	tsFrameTime
	tsLast
)

var font = text.NewAtlas(basicfont.Face7x13, text.ASCII)

type timeSeries struct {
	Values [tsLast][]int
	Text   [tsLast]*text.Text
}

func (t *timeSeries) append(entities, entityCap, memory, frameTime int) {
	t.Values[tsEntities] = append(t.Values[tsEntities], entities)
	t.Values[tsEntityCap] = append(t.Values[tsEntityCap], entityCap)
	t.Values[tsMemory] = append(t.Values[tsMemory], memory)
	t.Values[tsFrameTime] = append(t.Values[tsFrameTime], frameTime)
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
		text := text.New(px.V(0, 0), font)
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

// WorldStats visualizes statistics of the ECS world.
type WorldStats struct {
	HidePlots      bool
	HideArchetypes bool
	scale          float64
	drawer         imdraw.IMDraw
	tickRes        generic.Resource[resource.Tick]
	summary        *text.Text
	timeSeries     timeSeries
	frameTimer     frameTimer
	archetypes     archetypes
	text           *text.Text
}

// Initialize the system
func (s *WorldStats) Initialize(w *ecs.World, win *pixelgl.Window) {
	s.tickRes = generic.NewResource[resource.Tick](w)
	s.drawer = *imdraw.New(nil)

	s.scale = calcScaleCorrection()

	s.summary = text.New(px.V(0, 0), font)
	for i := 0; i < len(s.timeSeries.Text); i++ {
		s.timeSeries.Text[i] = text.New(px.V(0, 0), font)
	}
	fmt.Fprintf(s.timeSeries.Text[tsEntities], "Entities")
	fmt.Fprintf(s.timeSeries.Text[tsEntityCap], "Capacity")
	fmt.Fprintf(s.timeSeries.Text[tsMemory], "Memory")
	fmt.Fprintf(s.timeSeries.Text[tsFrameTime], "TPT")

	s.text = text.New(px.V(0, 0), font)
}

// Update the drawer.
func (s *WorldStats) Update(w *ecs.World) {
	if s.tickRes.Has() {
		tick := s.tickRes.Get().Tick
		s.frameTimer.Update(tick, time.Now())
	}
}

// Draw the system
func (s *WorldStats) Draw(w *ecs.World, win *pixelgl.Window) {
	var tick int64 = -1
	if s.tickRes.Has() {
		tick = s.tickRes.Get().Tick
	}
	stats := w.Stats()
	s.archetypes.Update(stats)

	s.summary.Clear()
	mem, units := toMemText(stats.Memory)
	fmt.Fprintf(
		s.summary, "Tick: %7d  |  Entities: %7d  |  Comp: %3d  |  Arch: %3d  |  Mem: %6.1f %s  |  TPS: %6.1f | TPT: %6.2f ms",
		tick, stats.Entities.Used, stats.ComponentCount, len(stats.Archetypes), mem, units, s.frameTimer.FPS(), float64(s.frameTimer.FrameTime().Microseconds())/1000,
	)

	numArch := len(stats.Archetypes)
	maxCapacity := 0
	for i := 0; i < numArch; i++ {
		cap := stats.Archetypes[i].Capacity
		if cap > maxCapacity {
			maxCapacity = cap
		}
	}
	dr := &s.drawer
	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()
	x0 := 10.0
	y0 := height - 20.0

	s.summary.Draw(win, px.IM.Moved(px.V(x0, y0)))
	y0 -= 10

	if !s.HidePlots {
		s.timeSeries.append(stats.Entities.Used, stats.Entities.Total, stats.Memory, int(s.frameTimer.FrameTime().Nanoseconds()))

		plotY0 := y0
		plotHeight := (height - 60) / 3
		if plotHeight > 150 {
			plotHeight = 150
		}
		plotWidth := (width - 20) * 0.25
		if s.HideArchetypes {
			plotWidth = width - 20
		}
		s.drawPlot(win, x0, plotY0-plotHeight, plotWidth, plotHeight, tsEntities, tsEntityCap)
		plotY0 -= plotHeight + 10
		s.drawPlot(win, x0, plotY0-plotHeight, plotWidth, plotHeight, tsMemory)
		plotY0 -= plotHeight + 10
		s.drawPlot(win, x0, plotY0-plotHeight, plotWidth, plotHeight, tsFrameTime)

		x0 += plotWidth + 10
	}

	archHeight := (y0 - 10) / float64(numArch)
	if !s.HideArchetypes {
		if archHeight >= 5 {
			archWidth := width - x0 - 10
			if archHeight > 30 {
				archHeight = 30
			}
			for i := 0; i < numArch; i++ {
				s.drawArchetype(
					win, x0, y0-float64(i+1)*archHeight, archWidth, archHeight,
					maxCapacity, &stats.Archetypes[i], s.archetypes.Components[i],
				)
			}
		} else {
			s.text.Clear()
			fmt.Fprintf(s.text, "Too many archetypes")
			s.text.Draw(win, px.IM.Moved(px.V(x0, y0-10)))
		}
	}

	dr.Draw(win)
	dr.Clear()
}

func (s *WorldStats) drawArchetype(win *pixelgl.Window, x, y, w, h float64, max int, arch *stats.ArchetypeStats, text *text.Text) {
	dr := &s.drawer

	cap := float64(arch.Capacity) / float64(max)
	cnt := float64(arch.Size) / float64(max)

	dr.Color = color.RGBA{160, 40, 40, 255}
	dr.Push(px.V(x, y), px.V(x+w*cnt, y+h))
	dr.Rectangle(0)
	dr.Reset()

	dr.Color = color.RGBA{40, 40, 160, 255}
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

func (s *WorldStats) drawPlot(win *pixelgl.Window, x, y, w, h float64, series ...timeSeriesType) {
	dr := &s.drawer

	dr.Color = color.RGBA{0, 25, 10, 255}
	dr.Push(px.V(x, y), px.V(x+w, y+h))
	dr.Rectangle(0)
	dr.Reset()

	yMax := 0
	for _, series := range series {
		values := s.timeSeries.Values[series]
		for _, v := range values {
			if v > yMax {
				yMax = v
			}
		}
	}

	dr.Color = color.White
	for _, series := range series {
		values := s.timeSeries.Values[series]
		yMax := 0
		for _, v := range values {
			if v > yMax {
				yMax = v
			}
		}

		numValues := len(values)
		if numValues > 0 {
			xStep := w / float64(len(values)-1)
			yScale := 0.95 * h / float64(yMax)

			for i := 0; i < numValues-1; i++ {
				xi := float64(i)
				x1, x2 := xi*xStep, xi*xStep+xStep
				y1, y2 := float64(values[i])*yScale, float64(values[i+1])*yScale

				dr.Push(px.V(x+x1, y+y1), px.V(x+x2, y+y2))
				dr.Line(1)
				dr.Reset()
			}
		}
	}

	dr.Push(px.V(x, y), px.V(x+w, y+h))
	dr.Rectangle(1)
	dr.Reset()

	dr.Draw(win)
	dr.Clear()

	if len(series) > 0 {
		text := s.timeSeries.Text[series[0]]
		text.Draw(win, px.IM.Moved(px.V(x+w-text.Bounds().W()-3, y+3)))
	}
}

func toMemText(bytes int) (float64, string) {
	if bytes < 10*1_048_576 {
		return float64(bytes) / 1024, "kB"
	}
	return float64(bytes) / 1_048_576, "MB"
}
