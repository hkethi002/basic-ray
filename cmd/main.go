package main

import (
	geometry "basic-ray/pkg/geometry"
	myio "basic-ray/pkg/io"
	render "basic-ray/pkg/render"
	_ "fmt"
)

func main() {
	eye := geometry.Point{0, 0, 0}
	bottomLeftCorner := geometry.Point{-2.5, -1.40625, -.5}
	bottomRightCorner := geometry.Point{2.5, -1.40625, -.5}
	topLeftCorner := geometry.Point{-2.5, 1.40625, -.5}
	camera := render.MakeCamera(bottomLeftCorner, bottomRightCorner, topLeftCorner, 2560, 1440)
	triangle1 := &geometry.Triangle{
		Vertex0:       geometry.Point{-5, -2, 1},
		Vertex1:       geometry.Point{-5, -2, -10},
		Vertex2:       geometry.Point{5, -2, -10},
		DiffuseAlbedo: [3]float64{.9, 0, 0},
		MaterialType:  2,
	}
	triangle2 := &geometry.Triangle{
		Vertex0:       geometry.Point{-5, -2, 1},
		Vertex1:       geometry.Point{5, -2, 1},
		Vertex2:       geometry.Point{5, -2, -10},
		DiffuseAlbedo: [3]float64{.9, 0, 0},
		MaterialType:  2,
	}
	triangleBlue := &geometry.Triangle{
		Vertex0:       geometry.Point{-5, -1.5, -5},
		Vertex1:       geometry.Point{-5, 2, -4.5},
		Vertex2:       geometry.Point{-5, -1.5, -4},
		DiffuseAlbedo: [3]float64{.1, .1, .9},
		MaterialType:  2,
	}

	lightSource := &render.LightSource{Location: geometry.Point{0, 5, 0}, RGB: render.Color{255, 255, 255}}
	triangles := []*geometry.Triangle{triangle1, triangle2, triangleBlue}

	render.Main(eye, lightSource, camera, triangles)

	myio.Write(camera, "output.ppm")

	// for i := 1.0; i < 9.0; i += .5 {
	// 	eye = geometry.Point{0, 0, 0 - i}
	// 	bottomLeftCorner = geometry.Point{-2.5, -1.40625, -.5 - i}
	// 	bottomRightCorner = geometry.Point{2.5, -1.40625, -.5 - i}
	// 	topLeftCorner = geometry.Point{-2.5, 1.40625, -.5 - i}
	// 	camera = render.MakeCamera(bottomLeftCorner, bottomRightCorner, topLeftCorner, 2560, 1440)
	// 	render.Main(eye, lightSource, camera, triangles)
	// 	myio.Write(camera, fmt.Sprintf("output_%f.ppm", i))
	// }
}
