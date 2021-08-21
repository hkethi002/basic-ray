package render

import (
	geometry "basic-ray/pkg/geometry"
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
