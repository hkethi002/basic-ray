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

	// Needed for area lighting
	GeometricFactor(shadeRec *ShadeRec) float64
	PDF(shadeRec *ShadeRec) float64
}

type PhysicalLightSource interface {
	GetDirection(shadeRec *ShadeRec) geometry.Vector
	IncidentRadiance(shadeRec *ShadeRec) Color
	CastsShadows() bool
	InShadow(ray *geometry.Ray, shadeRec *ShadeRec) bool

	// Needed for area lighting
	GeometricFactor(shadeRec *ShadeRec) float64
	PDF(shadeRec *ShadeRec) float64
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

func (lightSource *AmbientLight) GeometricFactor(shadeRec *ShadeRec) float64 {
	return 1.0
}

func (lightSource *AmbientLight) PDF(shadeRec *ShadeRec) float64 {
	return 1.0
}

type AmbientOccluder struct {
	Sampler      Sampler
	MinimumLight float64
	u, v, w      geometry.Vector
	BasicLight
}

func (lightSource *AmbientOccluder) GetDirection(shadeRec *ShadeRec) geometry.Vector {
	samplePoint := lightSource.Sampler.SampleHemisphere()
	return geometry.Normalize(geometry.Add(
		geometry.Add(
			geometry.ScalarProduct(lightSource.u, samplePoint[0]),
			geometry.ScalarProduct(lightSource.v, samplePoint[1]),
		),
		geometry.ScalarProduct(lightSource.w, samplePoint[2]),
	))
}

func (lightSource *AmbientOccluder) IncidentRadiance(shadeRec *ShadeRec) Color {
	lightSource.w = shadeRec.Normal
	// Avoid the vector pointing straight up
	lightSource.v = geometry.Normalize(geometry.CrossProduct(lightSource.w, geometry.Vector{0.0072, 1.0, 00.34}))
	lightSource.u = geometry.CrossProduct(lightSource.v, lightSource.w)

	shadowRay := geometry.Ray{Origin: shadeRec.HitPoint, Vector: lightSource.GetDirection(shadeRec)}

	if lightSource.InShadow(&shadowRay, shadeRec) {
		return ScalarProduct(
			lightSource.Color,
			lightSource.RadianceScalingFactor*lightSource.MinimumLight,
		)
	} else {
		return ScalarProduct(lightSource.Color, lightSource.RadianceScalingFactor)
	}
}

func (lightSource *AmbientOccluder) InShadow(ray *geometry.Ray, shadeRec *ShadeRec) bool {
	t := math.Inf(1)
	for _, object := range shadeRec.World.Objects {
		if object.ShadowHit(ray, &t) {
			return true
		}
	}
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

func (lightSource *DirectionalLight) PDF(shadeRec *ShadeRec) float64 {
	return 1.0
}

func (lightSource *DirectionalLight) GeometricFactor(shadeRec *ShadeRec) float64 {
	return 1.0
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

// pdf is probability density function
type LightObject interface {
	SampleSurface() (geometry.Point, geometry.Vector)
	PDF(shadeRec *ShadeRec) float64
}
