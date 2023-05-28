package cmd

import (
	geometry "basic-ray/pkg/geometry"
	lighting "basic-ray/pkg/lights"
	obj "basic-ray/pkg/objects"
	render "basic-ray/pkg/render"
	sceneIo "basic-ray/pkg/scene"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/spf13/cobra"
	"time"
	"fmt"
)

func init() {
	var output string
	var samples int
	var display bool

	var renderCmd = &cobra.Command{
		Use:   "render",
		Short: "Render a scene",
		// Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			RenderScene(output, samples, display)
		},
	}
	renderCmd.Flags().StringVarP(&output, "output", "o", "output.ppm", "output file name")
	renderCmd.Flags().IntVarP(&samples, "samples", "s", 1, "number of ray samples per pixel")
	renderCmd.Flags().BoolVarP(&display, "display", "d", false, "display rendering in progress")
	rootCmd.AddCommand(renderCmd)
}

func RenderScene(output string, samples int, display bool) {
	diskSampler := render.CreateJitteredSampler(samples, 83, 1)
	diskSampler.GenerateSamples()
	diskSampler.MapSamplesToCircle()

	objects := make([]render.GeometricObject, 6)
	objects[0] = &render.Sphere{Center: geometry.Point{0, -70, 0}, Radius: 10, Mesh: render.Mesh{KEpsilon: 0.001}}
	objects[0].(*render.Sphere).Material = &render.PhongMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.45,
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
	objects[0].(*render.Sphere).Shadows = true
	objects[1] = &render.Sphere{Center: geometry.Point{230, 30, 0}, Radius: 60, Mesh: render.Mesh{KEpsilon: 0.001}}
	objects[1].(*render.Sphere).Material = &render.PhongMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.45,
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
	objects[1].(*render.Sphere).Shadows = true
	objects[2] = &render.Plane{Point: geometry.Point{0, -150, 0}, Normal: geometry.Vector{0, 1, 0}, Mesh: render.Mesh{KEpsilon: 0.001}}
	objects[2].(*render.Plane).Material = &render.MatteMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.45,
			DiffuseColor:                 render.Color{1, 1, 1},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.65,
			DiffuseColor:                 render.Color{1, 1, 1},
		},
	}
	objects[2].(*render.Plane).Shadows = true
	objects[3] = obj.CreateDisk(geometry.Point{-400, 100, 300}, 100, geometry.Vector{1, 0, -1})

	emmisiveMaterial := &render.EmmisiveMaterial{
		RadianceScalingFactor: 10.0,
		Color:                 render.WHITE,
	}
	objects[3].(*obj.Disk).Material = emmisiveMaterial
	objects[3].(*obj.Disk).Sampler = diskSampler
	objects[3].(*obj.Disk).Shadows = false

	objects[4] = &render.Sphere{Center: geometry.Point{300, 150, 500}, Radius: 80, Mesh: render.Mesh{KEpsilon: 0.001}}
	objects[4].(*render.Sphere).Material = &render.MatteMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.45,
			DiffuseColor:                 render.Color{1, 0, 0},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.65,
			DiffuseColor:                 render.Color{1, 1, 1},
		},
	}
	objects[4].(*render.Sphere).Shadows = true
	objects[5] = &render.Sphere{Center: geometry.Point{100, -50, -190}, Radius: 80, Mesh: render.Mesh{KEpsilon: 0.001}}
	objects[5].(*render.Sphere).Material = &render.MatteMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.45,
			DiffuseColor:                 render.Color{1, 0.5, 0.5},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.65,
			DiffuseColor:                 render.Color{1, 0.5, 0.5},
		},
	}
	objects[5].(*render.Sphere).Shadows = true

	var sampler render.Sampler
	if samples == 1 {
		sampler = render.CreateRegularSampler(samples, 83, 1)
	} else {
		sampler = render.CreateJitteredSampler(samples, 83, 1)
	}
	viewPlane := render.ViewPlane{HorizontalResolution: 4096, VerticalResolution: 2160, PixelSize: 0.25, Gamma: 1, Sampler: sampler}
	// viewPlane := render.ViewPlane{HorizontalResolution: 800, VerticalResolution: 800, PixelSize: 0.5, Gamma: 1, Sampler: sampler}

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
	// 	DistanceToViewPlane: 300,
	// 	LookPoint:           geometry.Point{0, 25, 0},
	// 	Eye:                 geometry.Point{0, 25, -500},
	// 	UpVector:            geometry.Vector{0, 1, 0},
	// 	BaseCamera:          render.BaseCamera{ViewPlane: viewPlane, Pixels: &pixels},
	// 	Zoom:                1.0,
	// }
	camera.Display = true

	camera.Initialize()

	var lightSources []render.LightSource
	light := lighting.AreaLight{Object: objects[3].(*obj.Disk), Material: emmisiveMaterial}
	light.Shadows = true
	// light := render.PointLight{
	// 	Location:   geometry.Point{100, 250, -150},
	// 	BasicLight: render.BasicLight{Shadows: true, Color: render.WHITE, RadianceScalingFactor: 2.0}}
	// light2 := render.DirectionalLight{
	// 	Direction:  geometry.Vector{0, -1, 0},
	// 	BasicLight: render.BasicLight{Color: render.WHITE, RadianceScalingFactor: 1.0}}
	// light.Initialize()
	// light2.Initialize()
	lightSources = append(lightSources, &light)

	world := render.World{Camera: &camera, Lights: lightSources, Objects: objects[:4]}

	world.Shading = "area"
	ambientLight := render.AmbientLight{
		BasicLight: render.BasicLight{RadianceScalingFactor: 1.0},
		// MinimumLight: 0.0,
		// Sampler:      sampler,
	}
	ambientLight.Initialize()
	world.AmbientLight = &ambientLight

	if display {
		go render.MultiThreadedMain(&world, samples)
		pixelgl.Run(displayRender(&camera))
	} else {
		render.MultiThreadedMain(&world, samples)
	}
	sceneIo.WritePNG(&camera, output, 0.9)

}


func displayRender(camera render.Camera) func() {
	return func() {
		cfg := pixelgl.WindowConfig{
			Title:  "Pixel Rocks!",
			Bounds: pixel.R(0, 0, float64(camera.GetViewPlane().HorizontalResolution), float64(camera.GetViewPlane().VerticalResolution)),
			VSync: true,
		}
		win, err := pixelgl.NewWindow(cfg)
		if err != nil {
			panic(err)
		}
		var (
			frames = 0
			second = time.Tick(time.Second)
		)

		for !win.Closed() {
			win.Update()
			win.Canvas().SetPixels(camera.GetDisplayPixels())

			frames++
			select {
			case <-second:
				win.SetTitle(fmt.Sprintf("%s | FPS: %d ", cfg.Title, frames))
				frames = 0
			default:
			}
		}
	}
}


func checkValues(x []uint8) (uint8, int) {
	var v uint8
	j := -1
	for i := 0; i<len(x); i++ {
		if x[i] > v {
			v = x[i]
			j = i
		}
	}
	return v, j
}

