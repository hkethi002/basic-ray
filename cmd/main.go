package main

import (
	geometry "basic-ray/pkg/geometry"
	myio "basic-ray/pkg/io"
	render "basic-ray/pkg/render"
	"fmt"
)

func main() {
	eye := geometry.Point{0, 0, 0}
	bottomLeftCorner := geometry.Point{-2.5, -1.40625, -.5}
	bottomRightCorner := geometry.Point{2.5, -1.40625, -.5}
	topLeftCorner := geometry.Point{-2.5, 1.40625, -.5}
	camera := render.MakeCamera(bottomLeftCorner, bottomRightCorner, topLeftCorner, 256, 144)
	triangle := &geometry.Triangle{
		Vertex0:       geometry.Point{0, -2, 1},
		Vertex1:       geometry.Point{-5, -2, -10},
		Vertex2:       geometry.Point{5, -2, -10},
		DiffuseAlbedo: [3]float64{.9, 0, 0},
		MaterialType:  2,
	}
	lightSource := &render.LightSource{Location: geometry.Point{0, 5, 0}, RGB: render.Color{255, 255, 255}}
	triangles := []*geometry.Triangle{triangle}

	render.Main(eye, lightSource, camera, triangles)
	fmt.Println(len(*camera.Pixels))

	myio.Write(camera, "output.ppm")
}
