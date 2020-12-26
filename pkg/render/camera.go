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

func GetPoint(camera *Camera, pixel [2]float64) geometry.Point {
	v := geometry.Add(geometry.ScalarProduct(camera.unitX, pixel[0]), geometry.ScalarProduct(camera.unitY, pixel[0]))
	return geometry.Point{v[0] + camera.origin[0], v[1] + camera.origin[1], v[2] + camera.origin[2]}
}
