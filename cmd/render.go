package cmd

import (
	geometry "basic-ray/pkg/geometry"
	render "basic-ray/pkg/render"
	sceneIo "basic-ray/pkg/scene"
	"github.com/spf13/cobra"
)

func init() {
	var output string
	var samples int

	var renderCmd = &cobra.Command{
		Use:   "render",
		Short: "Render a scene",
		// Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			RenderScene(output, samples)
		},
	}
	renderCmd.Flags().StringVarP(&output, "output", "o", "output.ppm", "output file name")
	renderCmd.Flags().IntVarP(&samples, "samples", "s", 1, "number of ray samples per pixel")
	rootCmd.AddCommand(renderCmd)
}

/*
func main() {
	eye := geometry.Point{0, 0, 0}
	bottomLeftCorner := geometry.Point{-2.5, -1.40625, -2}
	bottomRightCorner := geometry.Point{2.5, -1.40625, -2}
	topLeftCorner := geometry.Point{-2.5, 1.40625, -2}
	camera := render.MakeCamera(bottomLeftCorner, bottomRightCorner, topLeftCorner, 1980, 1080)
	// camera := render.MakeCamera(bottomLeftCorner, bottomRightCorner, topLeftCorner, 1280, 720)
	// camera := render.MakeCamera(bottomLeftCorner, bottomRightCorner, topLeftCorner, 256, 144)

	// rgb := render.Color{.18, 0, .18}
	lightSource := &render.DirectionalLight{Direction: geometry.Vector{1.5, -1, 0}, RGB: render.Color{2000, 2000, 2000}}
	// lightSource := &render.DeltaLight{Location: geometry.Point{-2, 2, 0}, RGB: render.Color{10000, 10000, 10000}}
	// lightSource2 := &render.DeltaLight{Location: geometry.Point{2, 0.5, -1}, RGB: render.Color{5000, 5000, 5000}}
	object, err := sceneIo.ReadObject("NewSphere.json")
	check(err)

	triangles := geometry.TriangulateObject(object)

	// object, err = sceneIo.ReadObject("scene.json")
	object, err = sceneIo.ReadObject("scene2.json")
	check(err)
	triangles = append(triangles, geometry.TriangulateObject(object)...)
	for i, t := range triangles {
		t.Id = i
	}
	fmt.Println("Starting Ray Tracing...")

	// render.Main(eye, []render.LightSource{lightSource}, camera, triangles)
	render.MultiThreadedMain(eye, []render.LightSource{lightSource}, camera, triangles)

	fmt.Println("finished, writing image...")

	sceneIo.WriteImage(camera, "output.ppm")
}
*/
func RenderScene(output string, samples int) {
	objects := make([]render.GeometricObject, 2)
	objects[0] = &render.Sphere{Center: geometry.Point{0, -25, 0}, Radius: 80}
	objects[0].(*render.Sphere).Material.Color = render.Color{1, 0, 0}
	objects[1] = &render.Sphere{Center: geometry.Point{0, 30, 0}, Radius: 60}
	objects[1].(*render.Sphere).Material.Color = render.Color{1, 1, 0}
	// objects[2] = &render.Plane{Point: geometry.Point{0, 0, 0}, Normal: geometry.Vector{0, 1, -1}}
	// objects[2].(*render.Plane).Material.Color = render.Color{0, 0.3, 0}

	sampler := render.JitteredPointSampler{BaseSampler: render.BaseSampler{NumberOfSamples: 9, NumberOfSets: 10}}
	sampler.GenerateSamples()
	viewPlane := render.ViewPlane{HorizontalResolution: 400, VerticalResolution: 400, PixelSize: 1, Gamma: 1, Sampler: &sampler}

	pixels := make([][]render.Color, viewPlane.HorizontalResolution)
	for i := range pixels {
		pixels[i] = make([]render.Color, viewPlane.VerticalResolution)
	}
	camera := render.OrthoCamera{ViewPlane: viewPlane, TopLeft: geometry.Point{0, -25, -160}, Direction: geometry.Vector{0, 0, 1}, Pixels: &pixels}

	var lightSources []render.LightSource

	render.MultiThreadedMain(&camera, lightSources, objects, samples)

	sceneIo.WriteImage(&camera, output)
}
