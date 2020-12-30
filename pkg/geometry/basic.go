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

type Triangle struct {
	Id            int
	Vertex0       Point
	Vertex1       Point
	Vertex2       Point
	Normal        Vector
	VertexNormals []Vector
	// Albedo are [0, 1] per rgb color
	DiffuseAlbedo      [3]float64
	SpecularAlbedo     [3]float64
	TranslucenseAlbedo [3]float64
	// What angle the light changes at
	RefractionIndex float64

	MaterialType int
}

type TextureProperties struct {
	DiffuseAlbedo      [3]float64 `json:"diffuse"`
	SpecularAlbedo     [3]float64 `json:"specular"`
	TranslucenseAlbedo [3]float64 `json:"translucense"`
	// What angle the light changes at
	RefractionIndex float64 `json:"refraction_index"`
	MaterialType    int     `json:"material_type"`
}

type Object struct {
	Textures      []TextureProperties `json:"textures"`
	TextureMap    []int               `json:"texture_map"`
	Vertexes      []Point             `json:"vertexes"`
	Faces         [][]int             `json:"faces"`
	Normals       []Vector            `json:"normals"`
	VertexNormals []Vector            `json:"vertex_normals,omitempty"`
}
