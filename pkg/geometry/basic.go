package geometry

const (
	REFLECTIVE                = iota
	REFLECTIVE_AND_REFRACTIVE = iota
	DIFFUSE                   = iota
)

type Point [3]float64

type Vector [3]float64

type Ray struct {
	Vector Vector
	Origin Point
}

type Triangle struct {
	Vertex0 Point
	Vertex1 Point
	Vertex2 Point
	// Albedo are [0, 1] per rgb color
	DiffuseAlbedo      [3]float64
	SpecularAlbedo     [3]float64
	TranslucenseAlbedo [3]float64
	// What angle the light changes at
	RefractionIndex float64

	MaterialType int
}
