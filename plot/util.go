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
