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
	// camera := render.MakeCamera(bottomLeftCorner, bottomRightCorner, topLeftCorner, 2560, 1440)

	// rgb := render.Color{.18, 0, .18}
	// lightSource := &render.DirectionalLight{Direction: geometry.Vector{0, -1, -1}, RGB: render.Color{2000, 2000, 2000}}
	lightSource := &render.DeltaLight{Location: geometry.Point{-2, 2, 0}, RGB: render.Color{10000, 10000, 10000}}
	lightSource2 := &render.DeltaLight{Location: geometry.Point{2, 0.5, -1}, RGB: render.Color{5000, 5000, 5000}}
	object, err := myio.ReadObject("sphere2.json")
	check(err)

	triangles := geometry.TriangulateObject(object)

	// object, err = myio.ReadObject("scene.json")
	object, err = myio.ReadObject("scene2.json")
	check(err)
	triangles = append(triangles, geometry.TriangulateObject(object)...)
	for i, t := range triangles {
		t.Id = i
	}
	fmt.Println("Starting Ray Tracing...")

	// render.Main(eye, []render.LightSource{lightSource}, camera, triangles)
	render.MultiThreadedMain(eye, []render.LightSource{lightSource, lightSource2}, camera, triangles)

	fmt.Println("finished, writing image...")

	myio.Write(camera, "output.ppm")
}
