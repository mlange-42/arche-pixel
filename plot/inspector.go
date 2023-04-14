package plot

import (
	"fmt"
	"reflect"

	px "github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/mlange-42/arche-model/resource"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// Inspector drawer for inspecting entities.
//
// Shows information of the entity indicated by the SelectedEntity resource ([github.com/mlange-42/arche-model/resource.SelectedEntity]).
type Inspector struct {
	FieldNames  bool // Shows field names of nested structs in components.
	selectedRes generic.Resource[resource.SelectedEntity]
	text        *text.Text
}

// Initialize the system
func (m *Inspector) Initialize(w *ecs.World, win *pixelgl.Window) {
	m.selectedRes = generic.NewResource[resource.SelectedEntity](w)

	m.text = text.New(px.V(0, 0), font)
}

// Update the drawer.
func (m *Inspector) Update(w *ecs.World) {}

// UpdateInputs handles input events of the previous frame update.
func (m *Inspector) UpdateInputs(w *ecs.World, win *pixelgl.Window) {}

// Draw the system
func (m *Inspector) Draw(w *ecs.World, win *pixelgl.Window) {
	if !m.selectedRes.Has() {
		return
	}
	sel := m.selectedRes.Get().Selected
	if sel.IsZero() {
		return
	}

	height := win.Canvas().Bounds().H()
	x0 := 10.0
	y0 := height - 20.0

	m.text.Clear()
	fmt.Fprintf(m.text, "Entity %+v\n\n", sel)

	if !w.Alive(sel) {
		fmt.Fprint(m.text, "  dead entity")
		m.text.Draw(win, px.IM.Moved(px.V(x0, y0)))
		return
	}

	mask := w.Mask(sel)
	bits := mask.TotalBitsSet()

	fieldFormat := "    %-20s %-16s = %v\n"
	if m.FieldNames {
		fieldFormat = "    %-20s %-16s = %+v\n"
	}

	for i := 0; i < ecs.MaskTotalBits && bits > 0; i++ {
		id := ecs.ID(i)
		if mask.Get(id) {
			tp, _ := w.ComponentType(id)
			ptr := w.Get(sel, id)
			val := reflect.NewAt(tp, ptr).Elem()

			fmt.Fprintf(m.text, "  %s\n", tp.Name())
			for i := 0; i < val.NumField(); i++ {
				f := val.Field(i)
				fmt.Fprintf(m.text, fieldFormat,
					tp.Field(i).Name, f.Type(), f.Interface())
			}
			fmt.Fprint(m.text, "\n")
			bits--
		}
	}

	m.text.Draw(win, px.IM.Moved(px.V(x0, y0)))
}
