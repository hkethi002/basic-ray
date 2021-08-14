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
	GetMaterial() Material
}

type Mesh struct {
	Material Material
}

func (mesh *Mesh) GetMaterial() Material {
	return mesh.Material
}

type Plane struct {
	Point    geometry.Point
	Normal   geometry.Vector
	KEpsilon float64 // TODO: Figure out if this can be hard coded into the method
	Mesh
}

type Sphere struct {
	Center   geometry.Point
	Radius   float64
	KEpsilon float64 // TODO: Figure out if this can be hard coded into the method
	Mesh
}

type Triangle struct {
	Id            int
	Vertex0       geometry.Point
	Vertex1       geometry.Point
	Vertex2       geometry.Point
	Normal        geometry.Vector
	VertexNormals []geometry.Vector
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
	Vertexes      []geometry.Point    `json:"vertexes"`
	Faces         [][]int             `json:"faces"`
	Normals       []geometry.Vector   `json:"normals"`
	VertexNormals []geometry.Vector   `json:"vertex_normals,omitempty"`
}
