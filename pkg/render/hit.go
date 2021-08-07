package render

import (
	geometry "basic-ray/pkg/geometry"
	"math"
)

func (plane *Plane) Hit(ray *geometry.Ray, tmin *float64, shadeRec *ShadeRec) bool {
	cos := geometry.DotProduct(ray.Vector, plane.Normal)
	if cos == 0 {
		return false
	}
	t := geometry.DotProduct(geometry.CreateVector(plane.Point, ray.Origin), plane.Normal) / cos

	if t > plane.KEpsilon {
		*tmin = t
		shadeRec.Normal = plane.Normal
		shadeRec.LocalHitPoint = geometry.Translate(ray.Origin, geometry.ScalarProduct(ray.Vector, t))
		return true
	} else {
		return false
	}
}

func (sphere *Sphere) Hit(ray *geometry.Ray, tmin *float64, shadeRec *ShadeRec) bool {
	centerToOrigin := geometry.CreateVector(sphere.Center, ray.Origin)
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
	if t > sphere.KEpsilon {
		*tmin = t
		shadeRec.Normal = geometry.Normalize(centerToOrigin)
		shadeRec.LocalHitPoint = geometry.Translate(ray.Origin, geometry.ScalarProduct(ray.Vector, t))
		return true
	}

	t = (-b + e) / (2 * a)

	if t > sphere.KEpsilon {
		*tmin = t
		shadeRec.Normal = geometry.Normalize(centerToOrigin)
		shadeRec.LocalHitPoint = geometry.Translate(ray.Origin, geometry.ScalarProduct(ray.Vector, t))
		return true
	}
	return false
}