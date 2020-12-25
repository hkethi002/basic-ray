package geometry

import "fmt"

type Point [3]float64

type Vector interface {
	i() float64
	j() float64
	k() float64
	dotProduct(vector Vector) float64
	crossProduct(vector Vector) Vector
	normalize() Vector
	magnitude() float64
}

type Triangle struct {
	vertex0 Point
	vertex1 Point
	vertex2 Point
	// Albedo are [0, 1] per rgb color
	diffuseAlbedo      [3]float64
	specularAlbedo     [3]float64
	translucenseAlbedo [3]float64
	// What angle the light changes at
	refractionIndex float64
}
