package cmd

import (
	geometry "basic-ray/pkg/geometry"
	// lighting "basic-ray/pkg/lights"
	obj "basic-ray/pkg/objects"
	physics "basic-ray/pkg/physics"
	render "basic-ray/pkg/render"
	sceneIo "basic-ray/pkg/scene"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	var output string
	var samples int
	var ticks int

	var renderTimeCmd = &cobra.Command{
		Use:   "render-time",
		Short: "Render a scene over a time span",
		// Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			RenderSceneTicks(output, samples, ticks)
		},
	}
	renderTimeCmd.Flags().StringVarP(&output, "output", "o", "output.ppm", "output file name")
	renderTimeCmd.Flags().IntVarP(&samples, "samples", "s", 1, "number of ray samples per pixel")
	renderTimeCmd.Flags().IntVarP(&ticks, "ticks", "t", 5, "number of time ticks to render")
	rootCmd.AddCommand(renderTimeCmd)
}

func RenderSceneTicks(output string, samples, ticks int) {
	plane := &render.Plane{Point: geometry.Point{0, -150, 0}, Normal: geometry.Vector{0, 1, 0}, Mesh: render.Mesh{KEpsilon: 0.0001}}
	plane.Material = &render.MatteMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.15,
			DiffuseColor:                 render.Color{1, 1, 1},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.15,
			DiffuseColor:                 render.Color{1, 1, 1},
		},
	}
	plane.Shadows = true

	var sampler render.Sampler
	if samples == 1 {
		sampler = render.CreateRegularSampler(samples, 83, 1)
	} else {
		sampler = render.CreateJitteredSampler(samples, 83, 1)
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

	camera.Initialize()

	ball := obj.CreateSphere(geometry.Point{-400, 100, 300}, 100)
	ball.Material = &render.MatteMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.45,
			DiffuseColor:                 render.Color{1, 0, 0},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.65,
			DiffuseColor:                 render.Color{1, 0, 0},
		},
	}
	ball.Shadows = true

	var lightSources []render.LightSource
	light := render.DirectionalLight{
		Direction:  geometry.Vector{0, -1, 1},
		BasicLight: render.BasicLight{Color: render.WHITE, RadianceScalingFactor: 1.0}}
	light.Initialize()
	lightSources = append(lightSources, &light)

	objects := []render.GeometricObject{plane, ball}
	world := render.World{Camera: &camera, Lights: lightSources, Objects: objects}

	world.Shading = ""
	ambientLight := render.AmbientLight{
		BasicLight: render.BasicLight{RadianceScalingFactor: 0.01},
	}
	ambientLight.Initialize()
	world.AmbientLight = &ambientLight
	physBall := physics.SimpleObject{Location: ball.Center, Velocity: geometry.Vector{0.0, 1.0, 0.0}, Mesh: ball}
	physicalObjects := []physics.PhysicObject{&physBall}

	for tick := 0; tick < ticks; tick++ {
		for _, physicalObject := range physicalObjects {
			physicalObject.Simulate(1)
		}
		render.MultiThreadedMain(&world, samples)
		fmt.Println(ball.Center)
		tickOutput := fmt.Sprintf("%03d-%s.png", tick, output)
		sceneIo.WritePNG(&camera, tickOutput, 0.9)
	}

}
