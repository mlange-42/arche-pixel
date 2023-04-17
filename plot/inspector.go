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
func (i *Inspector) Initialize(w *ecs.World, win *pixelgl.Window) {
	i.selectedRes = generic.NewResource[resource.SelectedEntity](w)

	i.text = text.New(px.V(0, 0), defaultFont)
	i.helpText = text.New(px.V(0, 0), defaultFont)

	fmt.Fprint(i.helpText, "Toggle [f]ields, [t]ypes, [v]alues or [n]ames")
}

// Update the drawer.
func (i *Inspector) Update(w *ecs.World) {}

// UpdateInputs handles input events of the previous frame update.
func (i *Inspector) UpdateInputs(w *ecs.World, win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyF) {
		i.HideFields = !i.HideFields
		return
	}
	if win.JustPressed(pixelgl.KeyT) {
		i.HideTypes = !i.HideTypes
		return
	}
	if win.JustPressed(pixelgl.KeyV) {
		i.HideValues = !i.HideValues
		return
	}
	if win.JustPressed(pixelgl.KeyN) {
		i.HideNames = !i.HideNames
		return
	}
}

// Draw the system
func (i *Inspector) Draw(w *ecs.World, win *pixelgl.Window) {
	i.helpText.Draw(win, px.IM.Moved(px.V(10, 10)))

	if !i.selectedRes.Has() {
		return
	}
	sel := i.selectedRes.Get().Selected
	if sel.IsZero() {
		return
	}

	height := win.Canvas().Bounds().H()
	x0 := 10.0
	y0 := height - 20.0

	i.text.Clear()
	fmt.Fprintf(i.text, "Entity %+v\n\n", sel)

	if !w.Alive(sel) {
		fmt.Fprint(i.text, "  dead entity")
		i.text.Draw(win, px.IM.Moved(px.V(x0, y0)))
		return
	}

	mask := w.Mask(sel)
	bits := mask.TotalBitsSet()

	for j := 0; j < ecs.MaskTotalBits && bits > 0; j++ {
		id := ecs.ID(j)
		if mask.Get(id) {
			tp, _ := w.ComponentType(id)
			ptr := w.Get(sel, id)
			val := reflect.NewAt(tp, ptr).Elem()

			fmt.Fprintf(i.text, "  %s\n", tp.Name())

			if !i.HideFields {
				for k := 0; k < val.NumField(); k++ {
					i.printField(i.text, tp, tp.Field(k), val.Field(k))
				}
				fmt.Fprint(i.text, "\n")
			}
			bits--
		}
	}

	i.text.Draw(win, px.IM.Moved(px.V(x0, y0)))
}

func (i *Inspector) printField(w io.Writer, tp reflect.Type, field reflect.StructField, value reflect.Value) {
	fmt.Fprintf(w, "    %-20s ", field.Name)
	if !i.HideTypes {
		fmt.Fprintf(w, "    %-16s ", value.Type())
	}
	if !i.HideValues {
		if i.HideNames {
			fmt.Fprintf(w, "= %v", value.Interface())
		} else {
			fmt.Fprintf(w, "= %+v", value.Interface())
		}
	}
	fmt.Fprint(i.text, "\n")
}
