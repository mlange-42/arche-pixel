package plot

import (
	"math"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/vgimg"
)

var preferredTicks = []float64{1, 2, 5, 10}

func calcScaleCorrection() float64 {
	width := 100.0
	c := vgimg.New(vg.Points(width), vg.Points(width))
	img := c.Image()
	return width / float64(img.Bounds().Dx())
}

func calcTicksStep(max float64, desired int) float64 {
	steps := float64(desired)
	approxStep := float64(max) / (steps - 1)
	stepPower := math.Pow(10, -math.Floor(math.Log10(approxStep)))
	normalizedStep := approxStep * stepPower
	for _, s := range preferredTicks {
		if s >= normalizedStep {
			normalizedStep = s
			break
		}
	}
	return normalizedStep / stepPower
}

type ringBuffer[T any] struct {
	data  []T
	start int
}

func newRingBuffer[T any](cap int) ringBuffer[T] {
	return ringBuffer[T]{
		data:  make([]T, 0, cap),
		start: 0,
	}
}

func (r *ringBuffer[T]) Len() int {
	return len(r.data)
}

func (r *ringBuffer[T]) Cap() int {
	return cap(r.data)
}

func (r *ringBuffer[T]) Get(idx int) T {
	return r.data[(r.start+idx)%r.Cap()]
}

func (r *ringBuffer[T]) Add(elem T) {
	if cap(r.data) > len(r.data) {
		r.data = append(r.data, elem)
		return
	}
	r.data[r.start] = elem
	r.start = (r.start + 1) % r.Cap()
}