package render

import (
	geometry "basic-ray/pkg/geometry"
)

type ShadeRec struct {
	ObjectHit     bool
	Material      Material
	HitPoint      geometry.Point
	LocalHitPoint geometry.Point
	Ray           geometry.Ray
	depth         int
	direction     geometry.Vector
	Normal        geometry.Vector
	RGBColor      Color  // TODO: Remove this and use material pointers
	World         *World // What is this?
}

type GeometricObject interface {
	Hit(ray *geometry.Ray, tmin *float64, shadeRec *ShadeRec) bool
	ShadowHit(ray *geometry.Ray, tmin *float64) bool
	GetMaterial() Material
}

type Mesh struct {
	Material Material
	Shadows  bool
	Sampler  Sampler // Only used for area light geometries
	KEpsilon float64 // TODO: Figure out if this can be hard coded into the method
}

func (mesh *Mesh) GetMaterial() Material {
	return mesh.Material
}

type Plane struct {
	Point  geometry.Point
	Normal geometry.Vector
	Mesh
}

type Triangle struct {
	Vertex                   geometry.Point
	VectorA, VectorB, Normal geometry.Vector
	Mesh
}

type TriangleMesh struct {
	Triangles []*Triangle
	Mesh
}

type Sphere struct {
	Center geometry.Point
	Radius float64
	Mesh
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
	Vertexes      []geometry.Point    `json:"vertexes"`
	Faces         [][]int             `json:"faces"`
	Normals       []geometry.Vector   `json:"normals"`
	VertexNormals []geometry.Vector   `json:"vertex_normals,omitempty"`
}
