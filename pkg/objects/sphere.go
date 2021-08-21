package objects

import (
	geometry "basic-ray/pkg/geometry"
	render "basic-ray/pkg/render"
	"math"
)

type Sphere struct {
	Center geometry.Point
	Radius float64
	render.Mesh
}

func (sphere *Sphere) Hit(ray *geometry.Ray, tmin *float64, shadeRec *render.ShadeRec) bool {
	centerToOrigin := geometry.CreateVector(ray.Origin, sphere.Center)
	a := geometry.DotProduct(ray.Vector, ray.Vector)
	b := geometry.DotProduct(geometry.ScalarProduct(centerToOrigin, 2), ray.Vector)
	c := geometry.DotProduct(centerToOrigin, centerToOrigin) - (sphere.Radius * sphere.Radius)
	var e, t float64
	discriminant := (b * b) - 4*a*c
	if discriminant < 0 {
		return false
	}
	e = math.Sqrt(discriminant)
	t = (-b - e) / (2 * a)
	if t > sphere.KEpsilon && t < *tmin {
		*tmin = t
		shadeRec.Normal = geometry.ScalarProduct(geometry.Add(centerToOrigin, geometry.ScalarProduct(ray.Vector, t)), 1.0/sphere.Radius)
		shadeRec.LocalHitPoint = geometry.Translate(ray.Origin, geometry.ScalarProduct(ray.Vector, t))
		return true
	}

	t = (-b + e) / (2 * a)

	if t > sphere.KEpsilon && t < *tmin {
		*tmin = t
		shadeRec.Normal = geometry.ScalarProduct(geometry.Add(centerToOrigin, geometry.ScalarProduct(ray.Vector, t)), 1.0/sphere.Radius)
		shadeRec.LocalHitPoint = geometry.Translate(ray.Origin, geometry.ScalarProduct(ray.Vector, t))
		return true
	}
	return false
}

func (sphere *Sphere) ShadowHit(ray *geometry.Ray, tmin *float64) bool {
	if !sphere.Shadows {
		return false
	}
	centerToOrigin := geometry.CreateVector(ray.Origin, sphere.Center)
	a := geometry.DotProduct(ray.Vector, ray.Vector)
	b := geometry.DotProduct(geometry.ScalarProduct(centerToOrigin, 2), ray.Vector)
	c := geometry.DotProduct(centerToOrigin, centerToOrigin) - (sphere.Radius * sphere.Radius)
	var e, t float64
	discriminant := (b * b) - 4*a*c
	if discriminant < 0 {
		return false
	}
	e = math.Sqrt(discriminant)
	t = (-b - e) / (2 * a)
	if t > sphere.KEpsilon && t < *tmin {
		*tmin = t
		return true
	}

	t = (-b + e) / (2 * a)

	if t > sphere.KEpsilon && t < *tmin {
		*tmin = t
		return true
	}
	return false
}
