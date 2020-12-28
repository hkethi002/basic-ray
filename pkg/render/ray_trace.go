package render

import (
	geometry "basic-ray/pkg/geometry"
	_ "fmt"
	"math"
)

func Main(origin geometry.Point, lightSources []LightSource, camera *Camera, triangles []*geometry.Triangle) {
	var light Photon
	for i, row := range *camera.Pixels {
		for j, _ := range row {
			ray := geometry.Ray{
				Origin: origin,
				Vector: geometry.Normalize(geometry.CreateVector(GetPoint(camera, i, j), origin)),
			}
			light = Trace(&ray, triangles, lightSources, 0)
			(*camera.Pixels)[i][j] = light.rgb // GetWeightedColor()
		}
	}
}

func Trace(ray *geometry.Ray, triangles []*geometry.Triangle, lightSources []LightSource, depth int) Photon {
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

		photon = GetColor(ray, collision, triangle, lightSources, triangles, depth)
	}

	return photon
}

func GetColor(
	ray *geometry.Ray,
	reflectionPoint geometry.Point,
	triangle *geometry.Triangle,
	lightSources []LightSource,
	triangles []*geometry.Triangle,
	depth int,
) Photon {
	receiveVector := geometry.Normalize(geometry.ScalarProduct(ray.Vector, -1))
	switch triangle.MaterialType {
	case geometry.REFLECTIVE:
		reflectionRay := &geometry.Ray{Origin: reflectionPoint, Vector: GetReflectiveVector(ray.Vector, triangle)}
		return Trace(reflectionRay, triangles, lightSources, depth+1)
	case geometry.DIFFUSE:
		photons := GetDirectLight(reflectionPoint, triangles, lightSources)
		sampleRays := geometry.MakeSampleRays(reflectionPoint, triangle.GetNormal(), 16)
		var photon Photon
		for _, sampleRay := range sampleRays {
			photon = Trace(sampleRay, triangles, lightSources, depth+1)
			photons = append(photons, &photon)
		}

		return DiffuseShader(receiveVector, photons, triangle)

	}

	return Photon{}
}
