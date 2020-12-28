package main

import (
	geometry "basic-ray/pkg/geometry"
	myio "basic-ray/pkg/io"
	render "basic-ray/pkg/render"
	"fmt"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	eye := geometry.Point{0, 0, 0}
	bottomLeftCorner := geometry.Point{-2.5, -1.40625, -2}
	bottomRightCorner := geometry.Point{2.5, -1.40625, -2}
	topLeftCorner := geometry.Point{-2.5, 1.40625, -2}
	camera := render.MakeCamera(bottomLeftCorner, bottomRightCorner, topLeftCorner, 1280, 720)

	// rgb := render.Color{.18, 0, .18}
	lightSource := &render.DirectionalLight{Direction: geometry.Vector{5, -1, -1}, RGB: render.Color{1000, 1000, 1000}}
	object, err := myio.ReadObject("cube.json")
	check(err)
	// triangleBlue := &geometry.Triangle{
	// 	Vertex0:       geometry.Point{-5, -1.5, -5},
	// 	Vertex1:       geometry.Point{-5, 2, -4.5},
	// 	Vertex2:       geometry.Point{-5, -1.5, -4},
	// 	DiffuseAlbedo: [3]float64{.1, .1, .9},
	// 	MaterialType:  2,
	// }

	triangles := geometry.TriangulateObject(object)
	// triangles = append(triangles, triangleBlue)
	fmt.Println(triangles[2].Vertex0)
	fmt.Println(triangles[2].Vertex1)
	fmt.Println(triangles[2].Vertex2)

	render.Main(eye, lightSource, camera, triangles)

	myio.Write(camera, "output.ppm")

	// for i := 0.0; i < 9.0; i += .1 {
	// 	eye = geometry.Point{0, 0, 0 - i}
	// 	bottomLeftCorner = geometry.Point{-2.5, -1.40625, -.5 - i}
	// 	bottomRightCorner = geometry.Point{2.5, -1.40625, -.5 - i}
	// 	topLeftCorner = geometry.Point{-2.5, 1.40625, -.5 - i}
	// 	camera = render.MakeCamera(bottomLeftCorner, bottomRightCorner, topLeftCorner, 256, 144)
	// 	render.Main(eye, lightSource, camera, triangles)
	// 	myio.Write(camera, fmt.Sprintf("moving_down_red_road/output_%f.ppm", i))
	// }
}
