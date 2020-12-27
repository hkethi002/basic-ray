package render

import (
	geometry "basic-ray/pkg/geometry"
	"math"
)

func Main(origin geometry.Point, lightSource *LightSource, camera *Camera, triangles []*geometry.Triangle) {
	for i, row := range *camera.Pixels {
		for j, _ := range row {
			ray := geometry.Ray{
				Origin: origin,
				Vector: geometry.Normalize(geometry.CreateVector(GetPoint(camera, i, j), origin)),
			}
			(*camera.Pixels)[i][j] = Trace(&ray, triangles, lightSource, 0).rgb
		}
	}
}

func Trace(ray *geometry.Ray, triangles []*geometry.Triangle, lightSource *LightSource, depth int) Photon {
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
	lightSource *LightSource,
	triangles []*geometry.Triangle,
	depth int,
) Photon {
	receiveVector := geometry.ScalarProduct(ray.Vector, -1)
	switch triangle.MaterialType {
	case geometry.REFLECTIVE:
		reflectionRay := &geometry.Ray{Origin: reflectionPoint, Vector: GetReflectiveVector(ray.Vector, triangle)}
		return Trace(reflectionRay, triangles, lightSource, depth+1)
	case geometry.DIFFUSE:
		directLight := GetDirectLight(reflectionPoint, triangles, lightSource)
		if directLight != nil {
			return DiffuseShader(receiveVector, []*Photon{directLight}, triangle)
		}

	}

	return Photon{}
}
