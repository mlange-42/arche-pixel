package plot

import (
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/vgimg"
)

func calcScaleCorrection() float64 {
	width := 100.0
	c := vgimg.New(vg.Points(width), vg.Points(width))
	img := c.Image()
	return width / float64(img.Bounds().Dx())
}
