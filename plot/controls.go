package plot

import (
	"fmt"
	"image/color"
	"math"

	px "github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// Controls drawer and input handler.
//
// Allows to pause and resume a the simulation via a button or be pressing SPACE.
// Expects a world resource of type Systems ([github.com/mlange-42/arche-model/model/Systems])
type Controls struct {
	Scale      float64 // Spatial scaling: cell size in screen pixels. Optional, default 1.
	drawer     imdraw.IMDraw
	systemsRes generic.Resource[model.Systems]
	text       *text.Text
}

// Initialize the system
func (c *Controls) Initialize(w *ecs.World, win *pixelgl.Window) {
	if c.Scale <= 0 {
		c.Scale = 1
	}

	c.drawer = *imdraw.New(nil)
	c.text = text.New(px.V(0, 0), font)
}

// InitializeInputs initializes the InputHandler.
func (c *Controls) InitializeInputs(w *ecs.World, win *pixelgl.Window) {
	c.systemsRes = generic.NewResource[model.Systems](w)

	if !c.systemsRes.Has() {
		panic("resource of type Systems expected in Controls drawer")
	}
}

// Update the drawer.
func (c *Controls) Update(w *ecs.World) {}

// Inputs handles input events.
func (c *Controls) Inputs(w *ecs.World, win *pixelgl.Window) {
	sys := c.systemsRes.Get()
	if win.JustPressed(pixelgl.KeySpace) {
		sys.Paused = !sys.Paused
	}
	if win.JustPressed(pixelgl.MouseButton1) {
		width := win.Canvas().Bounds().W()
		height := win.Canvas().Bounds().H()
		pause := c.pauseBounds(width, height)

		mouse := win.MousePosition()
		if pause.Contains(mouse.X, mouse.Y) {
			sys.Paused = !sys.Paused
		}
	}
}

// Draw the system
func (c *Controls) Draw(w *ecs.World, win *pixelgl.Window) {
	dr := &c.drawer

	width := win.Canvas().Bounds().W()
	height := win.Canvas().Bounds().H()

	dr.Color = color.White

	pause := c.pauseBounds(width, height)

	dr.Push(px.V(pause.X, pause.Y), px.V(pause.X+pause.W, pause.Y+pause.H))
	dr.Rectangle(1)
	dr.Reset()

	dr.Draw(win)
	dr.Clear()

	c.text.Clear()

	sys := c.systemsRes.Get()
	if sys.Paused {
		fmt.Fprint(c.text, "Resume")
	} else {
		fmt.Fprint(c.text, "Pause")
	}

	wTxt := c.text.Bounds().W()
	hTxt := c.text.Bounds().H()
	cx, cy := pause.Center()
	c.text.Draw(win, px.IM.Scaled(px.V(wTxt/2, hTxt/2), c.Scale).Moved(px.V(math.Floor(cx-wTxt/2), math.Floor(cy-hTxt/2))))
}

func (c *Controls) pauseBounds(w, h float64) button {
	return button{
		w - 65*c.Scale,
		5,
		60 * c.Scale,
		30 * c.Scale,
	}
}

type button struct {
	X float64
	Y float64
	W float64
	H float64
}

func (b *button) Center() (float64, float64) {
	return b.X + b.W/2, b.Y + b.H/2
}

func (b *button) Contains(x, y float64) bool {
	return x >= b.X && y >= b.Y && x <= b.X+b.W && y <= b.Y+b.H
}
