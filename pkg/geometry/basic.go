package geometry

const (
	REFLECTIVE                = iota
	REFLECTIVE_AND_REFRACTIVE = iota
	FLAT_DIFFUSE              = iota
	GOURAUD_DIFFUSE           = iota
)

type Point [3]float64

type Vector [3]float64

type Ray struct {
	Vector Vector
	Origin Point
}
