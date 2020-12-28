package render

import (
	geometry "basic-ray/pkg/geometry"
	_ "fmt"
	"math"
)

func Main(origin geometry.Point, lightSource LightSource, camera *Camera, triangles []*geometry.Triangle) {
	var light Photon
	for i, row := range *camera.Pixels {
		for j, _ := range row {
			ray := geometry.Ray{
				Origin: origin,
				Vector: geometry.Normalize(geometry.CreateVector(GetPoint(camera, i, j), origin)),
			}
			light = Trace(&ray, triangles, lightSource, 0)
			(*camera.Pixels)[i][j] = light.rgb // GetWeightedColor()
		}
	}
}

func Trace(ray *geometry.Ray, triangles []*geometry.Triangle, lightSource LightSource, depth int) Photon {
	closestPoint := float64(math.Inf(1))
	var photon Photon
	if depth >= 4 {
		return photon
	}
	for _, triangle := range triangles {
		intersects := geometry.GetIntersection(ray, triangle)
		if intersects == nil {
			continue
		}
		collision := *intersects
		distance := geometry.Distance(collision, ray.Origin)

		if distance < closestPoint {
			closestPoint = distance
		} else {
			continue
		}

		photon = GetColor(ray, collision, triangle, lightSource, triangles, depth)
	}

	return photon
}

func GetColor(
	ray *geometry.Ray,
	reflectionPoint geometry.Point,
	triangle *geometry.Triangle,
	lightSource LightSource,
	triangles []*geometry.Triangle,
	depth int,
) Photon {
	receiveVector := geometry.Normalize(geometry.ScalarProduct(ray.Vector, -1))
	switch triangle.MaterialType {
	case geometry.REFLECTIVE:
		reflectionRay := &geometry.Ray{Origin: reflectionPoint, Vector: GetReflectiveVector(ray.Vector, triangle)}
		return Trace(reflectionRay, triangles, lightSource, depth+1)
	case geometry.DIFFUSE:
		directLight := GetDirectLight(reflectionPoint, triangles, lightSource)
		photons := make([]*Photon, 0)
		if directLight != nil {
			photons = append(photons, directLight)
		}
		sampleRays := geometry.MakeSampleRays(reflectionPoint, triangle.GetNormal(), 16)
		var photon Photon
		for _, sampleRay := range sampleRays {
			photon = Trace(sampleRay, triangles, lightSource, depth+1)
			photons = append(photons, &photon)
		}

		return DiffuseShader(receiveVector, photons, triangle)

	}

	return Photon{}
}
