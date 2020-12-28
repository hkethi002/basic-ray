package render

import (
	geometry "basic-ray/pkg/geometry"
	"math"
)

type Color [3]float64

type Photon struct {
	vector geometry.Vector
	rgb    Color
}

type LightSource interface {
	GetPhoton(destination geometry.Point) Photon
	GetDistance(destination geometry.Point) float64
}

type DeltaLight struct {
	Location geometry.Point
	RGB      Color
}

type DirectionalLight struct {
	Direction geometry.Vector
	RGB       Color
}

func (lightSource *DeltaLight) GetPhoton(destination geometry.Point) Photon {
	return Photon{vector: geometry.Normalize(geometry.CreateVector(destination, lightSource.Location)), rgb: lightSource.RGB}
}

func (lightSource *DeltaLight) GetDistance(destination geometry.Point) float64 {
	return geometry.Distance(destination, lightSource.Location)
}

func (lightSource *DirectionalLight) GetPhoton(destination geometry.Point) Photon {
	return Photon{vector: geometry.Normalize(lightSource.Direction), rgb: lightSource.RGB}
}

func (lightSource *DirectionalLight) GetDistance(destination geometry.Point) float64 {
	return float64(math.Inf(1))
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

func getDirectLightFromSingleSource(destination geometry.Point, triangles []*geometry.Triangle, lightSource LightSource) *Photon {
	lightDistance := lightSource.GetDistance(destination)
	photon := lightSource.GetPhoton(destination)
	ray := &geometry.Ray{Origin: destination, Vector: geometry.ScalarProduct(photon.vector, -1)}
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
	return &photon
}

func GetDirectLight(destination geometry.Point, triangles []*geometry.Triangle, lightSources []LightSource) []*Photon {
	photons := make([]*Photon, 0)
	var photon *Photon
	for _, lightSource := range lightSources {
		photon = getDirectLightFromSingleSource(destination, triangles, lightSource)
		if photon != nil {
			photons = append(photons, photon)
		}
	}
	return photons
}
