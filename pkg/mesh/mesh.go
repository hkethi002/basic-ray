package mesh

import (
	geometry "basic-ray/pkg/geometry"
)

type Mesh struct {
	Vertexes []geometry.Point
	Faces    [][]int
	Normals  []geometry.Vector
}

func CreateObject(
	mesh *Mesh,
	textures []geometry.TextureProperties,
	textureMap []int,
) *geometry.Object {
	if len(textures) == 0 {
		textures = append(textures, geometry.TextureProperties{DiffuseAlbedo: [3]float64{0.18, 0, 0.18}, MaterialType: 2})
	}

	if len(textureMap) == 0 {
		textureMap = make([]int, len(mesh.Faces))
	}

	return &geometry.Object{
		Vertexes:   mesh.Vertexes,
		Faces:      mesh.Faces,
		Normals:    mesh.Normals,
		Textures:   textures,
		TextureMap: textureMap,
	}
}
