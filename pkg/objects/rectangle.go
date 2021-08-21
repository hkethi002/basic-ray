package objects

import (
	geometry "basic-ray/pkg/geometry"
	render "basic-ray/pkg/render"
)

type Rect struct {
	Corner     geometry.Point
	A, B       geometry.Vector
	lenA, lenB float64
	render.Mesh

	Normal      geometry.Vector
	Sampler     render.Sampler
	inverseArea float64
}

func (rect *Rect) Hit(ray *geometry.Ray, tmin *float64, shadeRec *render.ShadeRec) bool {

	cos := geometry.DotProduct(ray.Vector, rect.Normal)
	if cos == 0 {
		return false
	}
	t := geometry.DotProduct(geometry.CreateVector(rect.Corner, ray.Origin), rect.Normal) / cos

	if t < rect.KEpsilon {
		return false
	}

	localHitPoint := geometry.Translate(ray.Origin, geometry.ScalarProduct(ray.Vector, t))

	vec := geometry.CreateVector(localHitPoint, rect.Corner)
	compA := geometry.DotProduct(vec, rect.A)
	compB := geometry.DotProduct(vec, rect.B)
	if compA < 0 || compA > (rect.lenA*rect.lenA) {
		return false
	}
	if compB < 0 || compB > (rect.lenB*rect.lenB) {
		return false
	}
	*tmin = t
	if cos > 0 {
		shadeRec.Normal = geometry.ScalarProduct(rect.Normal, -1)
	} else {
		shadeRec.Normal = rect.Normal
	}
	shadeRec.LocalHitPoint = localHitPoint
	return true
}

func (rect *Rect) ShadowHit(ray *geometry.Ray, tmin *float64) bool {
	cos := geometry.DotProduct(ray.Vector, rect.Normal)
	if cos == 0 {
		return false
	}

	// find where it hits the plane the rectangle is in
	t := geometry.DotProduct(geometry.CreateVector(rect.Corner, ray.Origin), rect.Normal) / cos

	if t < rect.KEpsilon {
		return false
	}

	localHitPoint := geometry.Translate(ray.Origin, geometry.ScalarProduct(ray.Vector, t))

	vec := geometry.CreateVector(localHitPoint, rect.Corner)
	compA := geometry.DotProduct(vec, rect.A)
	compB := geometry.DotProduct(vec, rect.B)
	if compA < 0 || compA > (rect.lenA*rect.lenA) {
		return false
	}
	if compB < 0 || compB > (rect.lenB*rect.lenB) {
		return false
	}

	*tmin = t
	return true
}

func CreateRect(corner geometry.Point, a, b geometry.Vector) *Rect {
	rect := Rect{
		Corner: corner,
		A:      a,
		B:      b,
		lenA:   geometry.Magnitude(a),
		lenB:   geometry.Magnitude(b),
		Normal: geometry.Normalize(geometry.CrossProduct(a, b)),
		Mesh:   render.Mesh{KEpsilon: 0.0001},
	}
	rect.inverseArea = 1 / (rect.lenA * rect.lenB)
	return &rect
}

func (rect *Rect) SampleSurface() (geometry.Point, geometry.Vector) {
	samplePoint := rect.Sampler.SampleUnitSquare()
	return geometry.Translate(
		rect.Corner,
		geometry.Add(
			geometry.ScalarProduct(rect.A, samplePoint[0]),
			geometry.ScalarProduct(rect.B, samplePoint[1]),
		),
	), rect.Normal
}

func (rect *Rect) PDF(shadeRec *render.ShadeRec) float64 {
	return rect.inverseArea
}
