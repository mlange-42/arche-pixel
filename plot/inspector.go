package plot

import (
	"fmt"
	"io"
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
// Entity selection is to be done by another system, e.g. by user input.
//
// Details can be adjusted using the HideXxx fields.
// Further, keys F, T, V and N can be used to toggle details during a running simulation.
type Inspector struct {
	HideFields  bool // Hides components fields.
	HideTypes   bool // Hides field types.
	HideValues  bool // Hides field values.
	HideNames   bool // Hide field names of nested structs.
	selectedRes generic.Resource[resource.SelectedEntity]
	text        *text.Text
	helpText    *text.Text
}

// Initialize the system
func (m *Inspector) Initialize(w *ecs.World, win *pixelgl.Window) {
	m.selectedRes = generic.NewResource[resource.SelectedEntity](w)

	m.text = text.New(px.V(0, 0), font)
	m.helpText = text.New(px.V(0, 0), font)

	fmt.Fprint(m.helpText, "Toggle [f]ields, [t]ypes, [v]alues or [n]ames")
}

// Update the drawer.
func (m *Inspector) Update(w *ecs.World) {}

// UpdateInputs handles input events of the previous frame update.
func (m *Inspector) UpdateInputs(w *ecs.World, win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyF) {
		m.HideFields = !m.HideFields
		return
	}
	if win.JustPressed(pixelgl.KeyT) {
		m.HideTypes = !m.HideTypes
		return
	}
	if win.JustPressed(pixelgl.KeyV) {
		m.HideValues = !m.HideValues
		return
	}
	if win.JustPressed(pixelgl.KeyN) {
		m.HideNames = !m.HideNames
		return
	}
}

// Draw the system
func (m *Inspector) Draw(w *ecs.World, win *pixelgl.Window) {
	m.helpText.Draw(win, px.IM.Moved(px.V(10, 10)))

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

	for i := 0; i < ecs.MaskTotalBits && bits > 0; i++ {
		id := ecs.ID(i)
		if mask.Get(id) {
			tp, _ := w.ComponentType(id)
			ptr := w.Get(sel, id)
			val := reflect.NewAt(tp, ptr).Elem()

			fmt.Fprintf(m.text, "  %s\n", tp.Name())

			if !m.HideFields {
				for i := 0; i < val.NumField(); i++ {
					m.printField(m.text, tp, tp.Field(i), val.Field(i))
				}
				fmt.Fprint(m.text, "\n")
			}
			bits--
		}
	}

	m.text.Draw(win, px.IM.Moved(px.V(x0, y0)))
}

func (m *Inspector) printField(w io.Writer, tp reflect.Type, field reflect.StructField, value reflect.Value) {
	fmt.Fprintf(w, "    %-20s ", field.Name)
	if !m.HideTypes {
		fmt.Fprintf(w, "    %-16s ", value.Type())
	}
	if !m.HideValues {
		if m.HideNames {
			fmt.Fprintf(w, "= %v", value.Interface())
		} else {
			fmt.Fprintf(w, "= %+v", value.Interface())
		}
	}
	fmt.Fprint(m.text, "\n")
}
