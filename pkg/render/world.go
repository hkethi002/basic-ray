package render

type World struct {
	Camera       Camera
	AmbientLight LightSource
	Lights       []LightSource
	Objects      []GeometricObject
	Shading      string
}
