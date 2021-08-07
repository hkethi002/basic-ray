package render

import (
	geometry "basic-ray/pkg/geometry"
)

const BIAS = 0.00001

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

/*
func (lightSource *DeltaLight) GetPhoton(destination geometry.Point) Photon {
	fallOff := math.Pow(lightSource.GetDistance(destination), 2)
	rgb := Color{
		lightSource.RGB[0] / fallOff,
		lightSource.RGB[1] / fallOff,
		lightSource.RGB[2] / fallOff,
	}

	return Photon{vector: geometry.Normalize(geometry.CreateVector(destination, lightSource.Location)), rgb: rgb}
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

func GetRefractedVector(incedentVector geometry.Vector, triangle *geometry.Triangle) geometry.Vector {
	normalVector := triangle.GetNormal()
	incedentVector = geometry.Normalize(incedentVector)
	c1 := geometry.DotProduct(normalVector, incedentVector)
	return geometry.Subtract(geometry.ScalarProduct(incedentVector, triangle.RefractionIndex), geometry.ScalarProduct(normalVector, (triangle.RefractionIndex*c1)))
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
		collisionDistance := geometry.Distance(destination, collision)
		if collisionDistance > lightDistance || collisionDistance < BIAS {
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
*/
