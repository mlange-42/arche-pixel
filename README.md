# Arche Pixel

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/arche-pixel/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/arche-pixel/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/arche-pixel)](https://goreportcard.com/report/github.com/mlange-42/arche-pixel)
[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/arche-pixel.svg)](https://pkg.go.dev/github.com/mlange-42/arche-pixel)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/arche-pixel)
[![MIT license](https://img.shields.io/github/license/mlange-42/arche-pixel)](https://github.com/mlange-42/arche-pixel/blob/main/LICENSE)

*Arche Pixel* provides OpenGL graphics and live plots for the [Arche](https://github.com/mlange-42/arche) Entity Component System (ECS) using the [Pixel](https://github.com/faiface/pixel) game engine.

<div align="center" width="100%">

![Screenshot](https://user-images.githubusercontent.com/44003176/232126308-60299642-0490-478d-82a5-48d862da6703.png)  
*Screenshot showing Arche Pixel features, visualizing an evolutionary forest model.*
</div>

## Features

* Free 2D drawing using a convenient OpenGL interface.
* Live plots using unified observers (time series, line, bar, scatter and contour plots).
* ECS engine monitor for detailed performance statistics.
* Entity inspector for debugging and inspection.
* Simulation controls to pause or limit speed interactively.
* User input handling for interactive simulations.

## Installation

```
go get github.com/mlange-42/arche-pixel
```

The dependencies of [go-gl/gl](https://github.com/go-gl/gl) and [go-gl/glfw](https://github.com/go-gl/glfw) apply. For Ubuntu/Debian-based systems, these are:

- `libgl1-mesa-dev`
- `xorg-dev`

## Usage

See the [API docs](https://pkg.go.dev/github.com/mlange-42/arche-pixel) for details and examples.

[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/arche-pixel.svg)](https://pkg.go.dev/github.com/mlange-42/arche-pixel)

## License

This project is distributed under the [MIT licence](./LICENSE).
