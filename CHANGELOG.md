# Changelog

## [[v0.7.0]](https://github.com/mlange-42/arche-pixel/compare/v0.6.0...v0.7.0)

### Breaking changes

* Upgrade to Arche v0.10.0 (#47, #48, #50)

## [[v0.6.0]](https://github.com/mlange-42/arche-pixel/compare/v0.5.1...v0.6.0)

### Features

* Adds a convenience function `window.Run(model)` for running a model on the main thread (#41)

### Bugfixes

* Fixes changed text alignment after migration to [gopxl/pixel v2](https://github.com/gopxl/pixel), in plots `Systems`, `Resources` and `Inspector` (#42)
* Destroy OpenGL windows on UI finalization (#45)

### Documentation

* Adds explanation for symbology and abbreviations to `plot.Monitor` (#43)

### Other

* Enable full tests with window creation using `xvfb` (#45)
* Add test coverage report to CI, add coveralls badge (#45)
* Add more tests for utility functions and different plots configurations (#46)

## [[v0.5.1]](https://github.com/mlange-42/arche-pixel/compare/v0.5.0...v0.5.1)

### Bugfixes

* Downgrade indirect dependencies from `gopxl/pixel` to fix `gopxl/mainthread` v2.1.0 crash on window creation (#40)

## [[v0.5.0]](https://github.com/mlange-42/arche-pixel/compare/v0.4.0...v0.5.0)

### Breaking changes

* Migrate to [gopxl/pixel v2](https://github.com/gopxl/pixel) (#39)
* Upgrade to Arche 0.9 (#39)
* Upgrade to Arche-Model 0.5 (#39)

### Infrastructure

* Upgrade to Go 1.21 toolchain (#39)

## [[v0.4.0]](https://github.com/mlange-42/arche-pixel/compare/v0.3.0...v0.4.0)

### Breaking changes

* Upgrade to Arche 0.8 (#35, #38)
* Upgrade to Arche-Model 0.4 (#38)

### Features

* `plot.Inspector` can be scrolled using arrow keys or mouse wheel (#36)
* New `plot.Systems` for inspecting ECS systems (#36)
* New `plot.Resources` for inspecting ECS resources (#36)
* `plot.Monitor` show number of cached filters (#37)
* `plot.Monitor` summary line wraps when window is not wide enough (#37)

## [[v0.3.0]](https://github.com/mlange-42/arche-pixel/compare/v0.2.0...v0.3.0)

### Breaking changes

* Drawer `plot.ImageRGB` uses one `MatrixLayers` observers instead of three `Matrix` observers (#31)
* Upgrade to `arche-model` v0.3.0 (#33)

### Features

* Drawer `plot.Lines` for plotting table observer data, with a line series per column, and a common X column (#22)
* Drawer `plot.Scatter` for plotting table observer data as scatter plots. Supports multiple observers and multiple series per observer (#25)
* Drawer `plot.Bars` for plotting row observer data as bar chart (#27)
* Drawer `plot.Contour` for plotting grid data as contours (#31)
* Drawer `plot.HeatMap` for plotting grid data as heat maps (#31)
* Drawer `plot.Field` for plotting 2D vector fields (#31)
* Plot title, axes labels and axes limits can be configured for plots (optional) (#30)
* Optional selection of columns in bar and time series plots (#30)
* Drawers `plot.ImageRGB` and `plot.Field` can freely assign layers to channels (#31)

### Bugfixes

* TimeSeries plot updates observer on every tick, not only every `UpdateInterval` ticks (#22)
* Plots that use `gonum/plot` don't crash on minimized window (#28)

### Other

* Plots use mono-spaced font and fixed tick label axis padding, to avoid jumping y axis (#26)
* Remove the last tick label from the x axis if close to the right margin, to avoid jumping x axis (#29)
* Scatter plots use solid instead of empty circle (#30)
* `Window` does not call `Drawer.Draw` when it is minimized (#32)

## [[v0.2.0]](https://github.com/mlange-42/arche-pixel/compare/v0.1.0...v0.2.0)

### Features

* New drawer `Inspector` for inspecting entities (#21)

## [[v0.1.0]](https://github.com/mlange-42/arche-pixel/compare/v0.0.3...v0.1.0)

### Features

* New drawer `PerfStats` for an overlay with performance stats in a corner of the window (#19)

### Other

* Upgrade to Arche v0.6.3 and Arche-Model v0.1.0 (#20)
* Promote to v0.1.0 to reflect increased API stability (#20)

## [[v0.0.3]](https://github.com/mlange-42/arche-pixel/compare/v0.0.2...v0.0.3)

### Breaking changes

* Renamed `Window.Add` to `Window.With`, taking drawer VarArgs and allows for chaining (#8, #11)
* `Drawer` interface has method `Update(w *ecs.World)` (#8)
* All plots are `Drawer` instead of `UISystem`, and are added to a `Window` (#8)
* Fields of `Bounds` renamed from `Width` and `Height` to `W` and `H` (#15)
* Upgrade to `arche-model` v0.0.5 (#16)

### Features

* Adds `Image` plot for plotting grids and matrices (#8)
* Adds `ImageRGB` plot for plotting multi-channel grids and matrices (#8)
* `Monitor` drawer for visualizing world and performance statistics (#10)
* Windows are resizable (#10)
* Window title can be set at construction time (#11)
* Adds method `UpdateInputs` to `Drawer` interface, for handling user input (#12, #14)
* Adds `Controls` plot and input handler for controlling simulation speed and pause via GUI or keyboard (#12)
* `Image` and `ImageRGB` auto-scale when no explicit scale is given (#13)
* Adds Method `window.Scale` to calculate scaling like in `Image` and `ImageRGB` (#13)

### Documentation

* Add separate examples for `Window` and `Drawer` (#9)

## [[v0.0.2]](https://github.com/mlange-42/arche-pixel/compare/v0.0.1...v0.0.2)

### Other

* Remove hard dependencies on resources `Tick` and `Termination` (#6)
* Upgrade dependency to Arche v0.6.1 and Arche-Model v0.0.2 (#6)