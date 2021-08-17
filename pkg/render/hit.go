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

	if t > plane.KEpsilon && t < *tmin {
		*tmin = t
		if cos > 0 {
			shadeRec.Normal = geometry.ScalarProduct(plane.Normal, -1)
		} else {
			shadeRec.Normal = plane.Normal
		}
		shadeRec.LocalHitPoint = geometry.Translate(ray.Origin, geometry.ScalarProduct(ray.Vector, t))
		return true
	} else {
		return false
	}
}

func (plane *Plane) ShadowHit(ray *geometry.Ray, tmin *float64) bool {
	if !plane.Shadows {
		return false
	}
	cos := geometry.DotProduct(ray.Vector, plane.Normal)
	if cos == 0 {
		return false
	}
	t := geometry.DotProduct(geometry.CreateVector(plane.Point, ray.Origin), plane.Normal) / cos

	if t > plane.KEpsilon && t < *tmin {
		*tmin = t
		return true
	}
	return false
}

func (sphere *Sphere) Hit(ray *geometry.Ray, tmin *float64, shadeRec *ShadeRec) bool {
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
