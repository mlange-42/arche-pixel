package window

import (
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/mlange-42/arche-model/model"
)

// Run is essential to run simulations that feature arche-pixel UI components.
// Call this function from the main function of your application, with a Model as argument.
// This is necessary, so that Model.Run runs on the main thread.
func Run(model *model.Model) {
	opengl.Run(model.Run)
}
