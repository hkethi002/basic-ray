package main

import (
	geometry "basic-ray/pkg/geometry"
	myio "basic-ray/pkg/io"
	mesh "basic-ray/pkg/mesh"
)

func main() {
	sphere := mesh.Sphere{
		Radius: 0.5,
		Origin: geometry.Point{0, -2, -5},
	}
	sphereMesh := sphere.CreateMesh(3, true)
	textures := []geometry.TextureProperties{
		geometry.TextureProperties{
			DiffuseAlbedo: [3]float64{0.18, 0, 0.18},
			MaterialType:  3,
		},
	}
	object := mesh.CreateObject(sphereMesh, textures, make([]int, 0))

	myio.WriteObject("sphere1.json", object)
}
