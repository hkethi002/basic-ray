package render

import (
	geometry "basic-ray/pkg/geometry"
)

type Camera struct {
	origin geometry.Point
	unitX  geometry.Vector // The normalized vector for the x axis in 3D space
	unitY  geometry.Vector // The normalized vector for the y axis in 3D space

	Pixels *[][]Color
}

func GetPoint(camera *Camera, pixelI int, pixelJ int) geometry.Point {
	v := geometry.Add(
		geometry.ScalarProduct(camera.unitX, float64(pixelI)),
		geometry.ScalarProduct(camera.unitY, float64(pixelJ)),
	)
	return geometry.Point{v[0] + camera.origin[0], v[1] + camera.origin[1], v[2] + camera.origin[2]}
}

func MakeCamera(bottomLeftCorner, bottomRightCorner, topLeftCorner geometry.Point, pixelWidth, pixelHeight int) *Camera {
	pixels := make([][]Color, pixelWidth)
	for i := range pixels {
		pixels[i] = make([]Color, pixelHeight)
	}
	return &Camera{
		origin: bottomLeftCorner,
		unitX:  geometry.ScalarProduct(geometry.CreateVector(bottomRightCorner, bottomLeftCorner), 1.0/float64(pixelWidth)),
		unitY:  geometry.ScalarProduct(geometry.CreateVector(topLeftCorner, bottomLeftCorner), 1.0/float64(pixelHeight)),
		Pixels: &pixels,
	}
}
