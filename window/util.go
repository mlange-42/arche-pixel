package window

import (
	"math"

	"github.com/faiface/pixel/pixelgl"
)

// CalcScale calculates the drawing scale for fitting a source region into a window's canvas.
func CalcScale(win *pixelgl.Window, srcWidth, srcHeight float64) float64 {
	winWidth, winHeight := win.Canvas().Bounds().W(), win.Canvas().Bounds().H()
	scX, scY := winWidth/float64(srcWidth), winHeight/float64(srcHeight)
	return math.Min(scX, scY)
}
