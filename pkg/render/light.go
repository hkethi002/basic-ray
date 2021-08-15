package render

import (
	geometry "basic-ray/pkg/geometry"
	"math"
)

const BIAS = 0.00001

type Color [3]float64

type Photon struct {
	vector geometry.Vector
	rgb    Color
}

type LightSource interface {
	GetDirection(shadeRec *ShadeRec) geometry.Vector
	IncidentRadiance(shadeRec *ShadeRec) Color
	CastsShadows() bool
	InShadow(ray *geometry.Ray, shadeRec *ShadeRec) bool
}

type BasicLight struct {
	Shadows               bool
	RadianceScalingFactor float64
	Color                 Color
}

func (lightSource *BasicLight) CastsShadows() bool {
	return lightSource.Shadows
}

type AmbientLight struct {
	BasicLight
}

func (lightSource *AmbientLight) GetDirection(shadeRec *ShadeRec) geometry.Vector {
	return geometry.Vector{0, 0, 0}
}

func (lightSource *AmbientLight) IncidentRadiance(shadeRec *ShadeRec) Color {
	return ScalarProduct(lightSource.Color, lightSource.RadianceScalingFactor)
}

func (lightSource *AmbientLight) InShadow(ray *geometry.Ray, shadeRec *ShadeRec) bool {
	return false
}

type AmbientOccluder struct {
	Sampler      Sampler
	MinimumLight Color
	BasicLight
}

// func (lightSource *AmbientOccluder) GetDirection(shadeRec *ShadeRec) geometry.Vector {
// 	return lightSource.Sampler.Sampler
// }

func (lightSource *AmbientOccluder) IncidentRadiance(shadeRec *ShadeRec) Color {
	return ScalarProduct(lightSource.Color, lightSource.RadianceScalingFactor)
}

func (lightSource *AmbientOccluder) InShadow(ray *geometry.Ray, shadeRec *ShadeRec) bool {
	return false
}

func (lightSource *BasicLight) Initialize() {
	if lightSource.Color == BLACK {
		lightSource.Color = WHITE
	}
	if lightSource.RadianceScalingFactor == 0 {
		lightSource.RadianceScalingFactor = 1.0
	}
}

type PointLight struct {
	Location geometry.Point
	BasicLight
}

func (lightSource *PointLight) GetDirection(shadeRec *ShadeRec) geometry.Vector {
	return geometry.Normalize(geometry.CreateVector(
		lightSource.Location,
		shadeRec.HitPoint,
	))
}

func (lightSource *PointLight) IncidentRadiance(shadeRec *ShadeRec) Color {
	return ScalarProduct(lightSource.Color, lightSource.RadianceScalingFactor)
}

func (lightSource *PointLight) InShadow(ray *geometry.Ray, shadeRec *ShadeRec) bool {
	d := geometry.Distance(lightSource.Location, ray.Origin)
	t := d

	for _, object := range shadeRec.World.Objects {
		if object.ShadowHit(ray, &t) && t < d {
			return true
		}
	}
	return false
}

type DirectionalLight struct {
	Direction geometry.Vector
	BasicLight
}

func (lightSource *DirectionalLight) GetDirection(shadeRec *ShadeRec) geometry.Vector {
	return geometry.ScalarProduct(lightSource.Direction, -1)
}

func (lightSource *DirectionalLight) IncidentRadiance(shadeRec *ShadeRec) Color {
	return ScalarProduct(lightSource.Color, lightSource.RadianceScalingFactor)
}

func (lightSource *DirectionalLight) InShadow(ray *geometry.Ray, shadeRec *ShadeRec) bool {
	var t float64
	d := math.Inf(1)

	for _, object := range shadeRec.World.Objects {
		if object.ShadowHit(ray, &t) && t < d {
			return true
		}
	}
	return false
}
