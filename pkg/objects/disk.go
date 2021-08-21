package objects

import (
	geometry "basic-ray/pkg/geometry"
	render "basic-ray/pkg/render"
)

type Disk struct {
	Center geometry.Point
	Radius float64
	Normal geometry.Vector
	render.Mesh

	Sampler     render.Sampler
	u, v        geometry.Vector // Create a 2D system
	inverseArea float64         // Avoid float64 division
}

func (disk *Disk) Hit(ray *geometry.Ray, tmin *float64, shadeRec *render.ShadeRec) bool {
	cos := geometry.DotProduct(ray.Vector, disk.Normal)
	if cos == 0 {
		return false
	}

	// find where it hits the plane the circle is in
	t := geometry.DotProduct(geometry.CreateVector(disk.Center, ray.Origin), disk.Normal) / cos
	if t < disk.KEpsilon {
		return false
	}

	localHitPoint := geometry.Translate(ray.Origin, geometry.ScalarProduct(ray.Vector, t))
	if geometry.Distance(localHitPoint, disk.Center) <= disk.Radius && t < *tmin {
		*tmin = t
		if cos > 0 {
			shadeRec.Normal = geometry.ScalarProduct(disk.Normal, -1)
		} else {
			shadeRec.Normal = disk.Normal
		}
		shadeRec.LocalHitPoint = localHitPoint
		return true
	} else {
		return false
	}
}

func (disk *Disk) ShadowHit(ray *geometry.Ray, tmin *float64) bool {
	cos := geometry.DotProduct(ray.Vector, disk.Normal)
	if cos == 0 {
		return false
	}

	// find where it hits the plane the circle is in
	t := geometry.DotProduct(geometry.CreateVector(disk.Center, ray.Origin), disk.Normal) / cos
	if t < disk.KEpsilon {
		return false
	}

	localHitPoint := geometry.Translate(ray.Origin, geometry.ScalarProduct(ray.Vector, t))
	if geometry.Distance(localHitPoint, disk.Center) <= disk.Radius && t < *tmin {
		*tmin = t
		return true
	} else {
		return false
	}
}

func (disk *Disk) SampleSurface() (geometry.Point, geometry.Vector) {
	samplePoint := disk.Sampler.SampleUnitCircle()
	return geometry.Translate(
		disk.Center,
		geometry.Add(
			geometry.ScalarProduct(disk.u, samplePoint[0]*disk.Radius),
			geometry.ScalarProduct(disk.v, samplePoint[1]*disk.Radius),
		),
	), disk.Normal
}

func (disk *Disk) PDF(shadeRec *render.ShadeRec) float64 {
	return disk.inverseArea
}

func CreateDisk(center geometry.Point, radius float64, normal geometry.Vector) *Disk {
	disk := Disk{Center: center, Radius: radius, Normal: normal, Mesh: render.Mesh{KEpsilon: 0.0001}}

	disk.inverseArea = 1.0 / (render.PI * radius * radius)
	disk.u = geometry.Normalize(geometry.CrossProduct(normal, geometry.Normalize(geometry.Vector{1, 100, -1})))
	disk.v = geometry.Normalize(geometry.CrossProduct(disk.u, normal))

	return &disk
}
