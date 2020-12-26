package render

import (
	geometry "basic-ray/pkg/geometry"
)

type Color [3]float64

type Photon struct {
	vector geometry.Vector
	rgb    Color
}

type LightSource struct {
	Location geometry.Point
	RGB      Color
}

func (lightSource *LightSource) GetPhoton(destination geometry.Point) Photon {
	return Photon{vector: geometry.Normalize(geometry.CreateVector(destination, lightSource.Location)), rgb: lightSource.RGB}
}

func GetReflectiveVector(incedentVector geometry.Vector, triangle *geometry.Triangle) geometry.Vector {
	normalVector := triangle.GetNormal()
	reflectionVector := geometry.Subtract(
		incedentVector,
		geometry.ScalarProduct(
			geometry.ScalarProduct(normalVector, geometry.DotProduct(incedentVector, normalVector)),
			2,
		),
	)
	return reflectionVector
}

func GetDirectLight(destination geometry.Point, triangles []*geometry.Triangle, lightSource *LightSource) *Photon {
	lightDistance := geometry.Distance(destination, lightSource.Location)
	ray := &geometry.Ray{Origin: destination, Vector: geometry.CreateVector(lightSource.Location, destination)}
	var photon Photon
	for _, triangle := range triangles {
		intersects := geometry.GetIntersection(ray, triangle)
		if intersects == nil {
			continue
		}
		collision := *intersects
		if geometry.Distance(destination, collision) > lightDistance {
			continue
		}
		return nil
	}
	photon = lightSource.GetPhoton(destination)
	return &photon
}
