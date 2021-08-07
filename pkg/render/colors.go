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

func ElementwiseProduct(colorA, colorB Color) Color {
	return Color{colorA[0] * colorB[0], colorA[1] * colorB[1], colorA[2] * colorB[2]}
}

func Pow(color Color, power float64) Color {
	return Color{math.Pow(color[0], power), math.Pow(color[1], power), math.Pow(color[2], power)}
}

var WHITE = Color{1, 1, 1}
var RED = Color{1, 0, 0}
var GREEN = Color{0, 1, 0}
var BLUE = Color{0, 0, 1}
var BLACK = Color{0, 0, 0}
