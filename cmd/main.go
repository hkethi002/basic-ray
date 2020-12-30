package cmd

import (
	geometry "basic-ray/pkg/geometry"
	render "basic-ray/pkg/render"
	sceneIo "basic-ray/pkg/scene"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	var output string

	var renderCmd = &cobra.Command{
		Use:   "render",
		Short: "Render a scene",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sceneFile := args[0]
			RenderScene(sceneFile, output)
		},
	}
	renderCmd.Flags().StringVarP(&output, "output", "o", "output.ppm", "output file name")
	rootCmd.AddCommand(renderCmd)
}

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

func RenderScene(sceneFile, output string) {
	fmt.Printf("Loading scene %s...\n", sceneFile)
	scene, err := sceneIo.LoadScene(sceneFile)
	check(err)

	camera, objects, lightSources, err := scene.GetComponents()
	check(err)

	triangles := make([]*geometry.Triangle, 0)
	for _, object := range objects {
		triangles = append(triangles, geometry.TriangulateObject(object)...)
	}

	fmt.Println("Starting Ray Tracing...")
	render.MultiThreadedMain(scene.Camera.Eye, lightSources, camera, triangles)

	fmt.Println("finished, writing image...")
	sceneIo.WriteImage(camera, output)
}
