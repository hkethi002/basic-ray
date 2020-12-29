package main

import (
	geometry "basic-ray/pkg/geometry"
	myio "basic-ray/pkg/io"
	mesh "basic-ray/pkg/mesh"
)

func main() {
	sphere := mesh.Sphere{
		Radius: 0.5,
		Origin: geometry.Point{0, 0, -3},
	}
	sphereMesh := sphere.CreateMesh(0)
	object := mesh.CreateObject(sphereMesh, make([]geometry.TextureProperties, 0), make([]int, 0))

	myio.WriteObject("sphere2.json", object)
}
