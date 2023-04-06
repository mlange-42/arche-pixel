# Changelog

## [[v0.0.3]](https://github.com/mlange-42/arche-pixel/compare/v0.0.2...v0.0.3)

### Breaking changes

* Renamed `Window.Add` to `Window.AddDrawer` (#8)
* `Drawer` interface has method `Update(w *ecs.World)` (#8)
* All plots are `Drawer` instead of `UISystem`, and are added to a `Window` (#8)

### Features

* Adds `Image` plot for plotting grids and matrices (#8)
* Adds `ImageRGB` plot for plotting multi-channel grids and matrices (#8)

## [[v0.0.2]](https://github.com/mlange-42/arche-pixel/compare/v0.0.1...v0.0.2)

### Other

* Remove hard dependencies on resources `Tick` and `Termination` (#6)
* Upgrade dependency to Arche v0.6.1 and Arche-Model v0.0.2 (#6)