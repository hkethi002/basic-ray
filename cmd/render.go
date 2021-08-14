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
	objects := make([]render.GeometricObject, 6)
	objects[0] = &render.Sphere{Center: geometry.Point{0, 25, 0}, Radius: 80}
	objects[0].(*render.Sphere).Material = &render.PhongMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.15,
			DiffuseColor:                 render.Color{1, 0, 1},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.65,
			DiffuseColor:                 render.Color{1, 0, 1},
		},
		GlossyBRDF: &render.GlossySpecular{
			SpecularReflectionCoefficient: 0.25,
			SpecularColor:                 render.Color{1, 0, 1},
			Exp:                           10,
		},
	}
	objects[1] = &render.Sphere{Center: geometry.Point{230, 30, 0}, Radius: 60}
	objects[1].(*render.Sphere).Material = &render.PhongMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.25,
			DiffuseColor:                 render.Color{1, 1, 0},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.6,
			DiffuseColor:                 render.Color{1, 1, 0},
		},
		GlossyBRDF: &render.GlossySpecular{
			SpecularReflectionCoefficient: 0.2,
			SpecularColor:                 render.Color{1, 1, 1},
			Exp:                           100,
		},
	}
	objects[2] = &render.Plane{Point: geometry.Point{0, -150, 0}, Normal: geometry.Vector{0, 1, 0}}
	objects[2].(*render.Plane).Material = &render.MatteMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.15,
			DiffuseColor:                 render.Color{.5, .5, .5},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.45,
			DiffuseColor:                 render.Color{.5, .5, .5},
		},
	}
	objects[3] = &render.Sphere{Center: geometry.Point{-400, -25, 500}, Radius: 80}
	objects[3].(*render.Sphere).Material = &render.MatteMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.25,
			DiffuseColor:                 render.Color{1, 0, 0},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.65,
			DiffuseColor:                 render.Color{1, 0, 0},
		},
	}
	objects[4] = &render.Sphere{Center: geometry.Point{300, 150, 500}, Radius: 80}
	objects[4].(*render.Sphere).Material = &render.MatteMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.25,
			DiffuseColor:                 render.Color{1, 1, 1},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.65,
			DiffuseColor:                 render.Color{1, 1, 1},
		},
	}
	objects[5] = &render.Sphere{Center: geometry.Point{100, -50, -190}, Radius: 80}
	objects[5].(*render.Sphere).Material = &render.MatteMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.25,
			DiffuseColor:                 render.Color{1, 0.5, 0.5},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.65,
			DiffuseColor:                 render.Color{1, 0.5, 0.5},
		},
	}

	var sampler render.Sampler
	if samples == 1 {
		sampler = &render.RegularSampler{BaseSampler: render.BaseSampler{NumberOfSamples: 1, NumberOfSets: 1}}
		sampler.(*render.RegularSampler).GenerateSamples()
	} else {
		sampler = render.CreateJitteredSampler(samples, 83, 0)
	}
	// viewPlane := render.ViewPlane{HorizontalResolution: 4096, VerticalResolution: 2160, PixelSize: 0.25, Gamma: 1, Sampler: sampler}
	viewPlane := render.ViewPlane{HorizontalResolution: 800, VerticalResolution: 800, PixelSize: 0.5, Gamma: 1, Sampler: sampler}

	pixels := make([][]render.Color, viewPlane.HorizontalResolution)
	for i := range pixels {
		pixels[i] = make([]render.Color, viewPlane.VerticalResolution)
	}
	camera := render.ThinLensCamera{
		DistanceToViewPlane: 300,
		LookPoint:           geometry.Point{0, 25, 0},
		Eye:                 geometry.Point{0, 25, -500},
		UpVector:            geometry.Vector{0, 1, 0},
		BaseCamera:          render.BaseCamera{ViewPlane: viewPlane, Pixels: &pixels},
		FocalDistance:       500,
		Zoom:                1,
		LensRadius:          10,
		Sampler:             sampler,
	}
	// camera := render.PinholeCamera{
	// 	DistanceToViewPlane: 650.0,
	// 	LookPoint:           geometry.Point{-5, 0, 0},
	// 	Eye:                 geometry.Point{0, 0, 500},
	// 	UpVector:            geometry.Vector{0, 1, 0},
	// 	BaseCamera:          render.BaseCamera{ViewPlane: viewPlane, Pixels: &pixels},
	// }

	camera.Initialize()

	var lightSources []render.LightSource
	light := render.PointLight{
		Location:   geometry.Point{100, 50, -150},
		BasicLight: render.BasicLight{Color: render.WHITE, RadianceScalingFactor: 2.0}}
	// light := render.DirectionalLight{
	// 	Direction:  geometry.Vector{0, -1, 0},
	// 	BasicLight: render.BasicLight{Color: render.WHITE, RadianceScalingFactor: 3.0}}
	light.Initialize()
	lightSources = append(lightSources, &light)

	world := render.World{Camera: &camera, Lights: lightSources, Objects: objects}

	ambientLight := render.AmbientLight{
		BasicLight: render.BasicLight{RadianceScalingFactor: 1.0},
	}
	ambientLight.Initialize()
	world.AmbientLight = &ambientLight

	render.MultiThreadedMain(&world, samples)

	sceneIo.WritePNG(&camera, output, 0.9)
}
