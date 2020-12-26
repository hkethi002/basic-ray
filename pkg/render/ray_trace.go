package render

import (
	geometry "basic-ray/pkg/geometry"
	"math"
)

func DoThing(ray *geometry.Ray, triangles []*geometry.Triangle) Color {
	closestPoint := float64(math.Inf(1))
	var closestColor Color
	for _, triangle := range triangles {
		intersects := geometry.GetIntersection(ray, triangle)
		if intersects == nil {
			continue
		}
		collision := *intersects

		if collision[2] < closestPoint {
			closestPoint = collision[2]
		} else {
			continue
		}

		closestColor = GetColor(collision, triangle)
	}
	return closestColor
}

func GetColor(reflectionPoint geometry.Point, triangle *geometry.Triangle) Color {
	return Color{}
}
