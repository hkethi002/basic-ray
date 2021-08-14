package render

import (
	"math"
)

func Add(colorA, colorB Color) Color {
	return Color{colorA[0] + colorB[0], colorA[1] + colorB[1], colorA[2] + colorB[2]}
}

func ScalarProduct(color Color, scalar float64) Color {
	return Color{color[0] * scalar, color[1] * scalar, color[2] * scalar}
}

func ElementwiseProduct(colors ...Color) Color {
	if len(colors) < 1 {
		return Color{}
	}
	color := Color{colors[0][0], colors[0][1], colors[0][2]}
	if len(colors) < 2 {
		return color
	}
	for c := 1; c < len(colors); c++ {
		color[0] = color[0] * colors[c][0]
		color[1] = color[1] * colors[c][1]
		color[1] = color[1] * colors[c][2]
	}
	return color
}

func Pow(color Color, power float64) Color {
	return Color{
		math.Pow(color[0], power),
		math.Pow(color[1], power),
		math.Pow(color[2], power),
	}
}

var WHITE = Color{1, 1, 1}
var RED = Color{1, 0, 0}
var GREEN = Color{0, 1, 0}
var BLUE = Color{0, 0, 1}
var BLACK = Color{0, 0, 0}
