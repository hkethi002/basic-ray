package lights

import (
	geometry "basic-ray/pkg/geometry"
	render "basic-ray/pkg/render"
	"math"
)

type AreaLight struct {
	lightNormal    geometry.Vector
	samplePoint    geometry.Point
	incidentVector geometry.Vector
	Object         render.LightObject
	Material       *render.EmmisiveMaterial
	render.BasicLight
}

func (light *AreaLight) GetDirection(shadeRec *render.ShadeRec) geometry.Vector {
	light.samplePoint, light.lightNormal = light.Object.SampleSurface()
	light.incidentVector = geometry.Normalize(geometry.CreateVector(light.samplePoint, shadeRec.HitPoint))
	return light.incidentVector
}

func (light *AreaLight) InShadow(ray *geometry.Ray, shadeRec *render.ShadeRec) bool {
	t := geometry.Distance(light.samplePoint, ray.Origin)

	for _, object := range shadeRec.World.Objects {
		if object.ShadowHit(ray, &t) {
			return true
		}
	}
	return false
}

func (light *AreaLight) IncidentRadiance(shadeRec *render.ShadeRec) render.Color {
	incidentCos := geometry.DotProduct(geometry.ScalarProduct(light.lightNormal, -1), light.incidentVector)

	if incidentCos > 0.0 {
		return light.Material.GetEmmittedRadiance(shadeRec)
	} else {
		return render.BLACK
	}
}

func (light *AreaLight) GeometricFactor(shadeRec *render.ShadeRec) float64 {
	incidentCos := geometry.DotProduct(geometry.ScalarProduct(light.lightNormal, -1), light.incidentVector)

	d2 := math.Pow(geometry.Distance(light.samplePoint, shadeRec.HitPoint), 2)

	return incidentCos / d2
}

func (light *AreaLight) PDF(shadeRec *render.ShadeRec) float64 {
	return light.Object.PDF(shadeRec)
}
