package geometry

import (
	_ "fmt"
	"math/rand"
)

func MakeSampleRays(origin Point, normalVector Vector, rayCount int) []*Ray {
	var ray Ray
	var v Vector
	sampleRays := make([]*Ray, 0)
	for i := 0; i < rayCount; i++ {
		v = Vector{
			float64(rand.Intn(10)),
			float64(rand.Intn(10)),
			float64(rand.Intn(10)),
		}
		if v[0] == 0 && v[1] == 0 && v[2] == 0 {
			v[rand.Intn(3)] += 1
		}
		ray = Ray{
			Vector: Normalize(v),
			Origin: origin,
		}
		if DotProduct(ray.Vector, normalVector) <= 0 {
			ray.Vector = ScalarProduct(ray.Vector, -1)
		}
		sampleRays = append(sampleRays, &ray)

	}
	return sampleRays
}
