# Changelog

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