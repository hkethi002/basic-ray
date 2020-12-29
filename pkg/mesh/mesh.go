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
	return &geometry.Object{
		Vertexes:   mesh.Vertexes,
		Faces:      mesh.Faces,
		Normals:    mesh.Normals,
		Textures:   textures,
		TextureMap: textureMap,
	}
}
