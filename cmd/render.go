package cmd

import (
	geometry "basic-ray/pkg/geometry"
	lighting "basic-ray/pkg/lights"
	obj "basic-ray/pkg/objects"
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

func RenderScene(output string, samples int) {
	// Disk for light
	diskSampler := render.CreateJitteredSampler(samples, 83, 1)
	emmisiveMaterial := &render.EmmisiveMaterial{
		RadianceScalingFactor: 5.0,
		Color:                 render.WHITE,
	}
	// disk := obj.CreateRect(geometry.Point{-100, 0, 100}, geometry.Vector{70, 0, 70}, geometry.Vector{0, 100, 0})
	disk := obj.CreateDisk(geometry.Point{-400, 100, 300}, 100, geometry.Vector{1, 0, -1})
	disk.Material = emmisiveMaterial
	disk.Sampler = diskSampler
	disk.Shadows = false

	// box thin
	topRect := obj.CreateRect(geometry.Point{0, -50, 0}, geometry.Vector{50, 0, 0}, geometry.Vector{0, 0, 50})
	topRect.Material = &render.MatteMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.45,
			DiffuseColor:                 render.Color{1, 1, 0},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.65,
			DiffuseColor:                 render.Color{1, 1, 0},
		},
	}

	rect1 := obj.CreateRect(geometry.Point{0, -150, 0}, geometry.Vector{0, 100, 0}, geometry.Vector{0, 0, 50})
	rect1.Material = &render.MatteMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.45,
			DiffuseColor:                 render.Color{1, 1, 0},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.65,
			DiffuseColor:                 render.Color{1, 1, 0},
		},
	}
	rect2 := obj.CreateRect(geometry.Point{0, -150, 50}, geometry.Vector{0, 100, 0}, geometry.Vector{50, 0, 0})
	rect2.Material = &render.MatteMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.45,
			DiffuseColor:                 render.Color{1, 1, 0},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.65,
			DiffuseColor:                 render.Color{1, 1, 0},
		},
	}
	rect3 := obj.CreateRect(geometry.Point{50, -150, 50}, geometry.Vector{0, 100, 0}, geometry.Vector{0, 0, -50})
	rect3.Material = &render.MatteMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.45,
			DiffuseColor:                 render.Color{1, 1, 0},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.65,
			DiffuseColor:                 render.Color{1, 1, 0},
		},
	}
	rect4 := obj.CreateRect(geometry.Point{50, -150, 0}, geometry.Vector{0, 100, 0}, geometry.Vector{-50, 0, 0})
	rect4.Material = &render.MatteMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.45,
			DiffuseColor:                 render.Color{1, 1, 0},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.65,
			DiffuseColor:                 render.Color{1, 1, 0},
		},
	}
	plane := &render.Plane{Point: geometry.Point{0, -150, 0}, Normal: geometry.Vector{0, 1, 0}, Mesh: render.Mesh{KEpsilon: 0.001}}
	plane.Material = &render.MatteMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.45,
			DiffuseColor:                 render.Color{1, 1, 1},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.65,
			DiffuseColor:                 render.Color{1, 1, 1},
		},
	}
	plane.Shadows = true

	sphere1 := &obj.Sphere{Center: geometry.Point{-300, 0, -200}, Radius: 100}
	sphere1.Material = &render.PhongMaterial{
		AmbientBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.25,
			DiffuseColor:                 render.Color{3, 0, 3},
		},
		DiffuseBRDF: &render.LambertianShader{
			DiffuseReflectionCoefficient: 0.45,
			DiffuseColor:                 render.Color{3, 0, 3},
		},
		GlossyBRDF: &render.GlossySpecular{
			SpecularReflectionCoefficient: 0.75,
			SpecularColor:                 render.Color{2, 0.5, 2},
			Exp:                           10,
		},
	}
	sphere1.Shadows = true

	objects := []render.GeometricObject{disk, rect1, rect2, rect3, rect4, plane, topRect, sphere1}

	var sampler render.Sampler
	if samples == 1 {
		sampler = render.CreateRegularSampler(samples, 83, 1)
	} else {
		sampler = render.CreateJitteredSampler(samples, 83, 1)
	}
	// viewPlane := render.ViewPlane{HorizontalResolution: 4096, VerticalResolution: 2160, PixelSize: 0.25, Gamma: 1, Sampler: sampler}
	viewPlane := render.ViewPlane{HorizontalResolution: 3200, VerticalResolution: 1600, PixelSize: 0.25, Gamma: 1, Sampler: sampler}

	pixels := make([][]render.Color, viewPlane.HorizontalResolution)
	for i := range pixels {
		pixels[i] = make([]render.Color, viewPlane.VerticalResolution)
	}
	camera := render.ThinLensCamera{
		DistanceToViewPlane: 200,
		LookPoint:           geometry.Point{0, 25, -100},
		Eye:                 geometry.Point{0, 25, -600},
		UpVector:            geometry.Vector{0, 1, 0},
		BaseCamera:          render.BaseCamera{ViewPlane: viewPlane, Pixels: &pixels},
		FocalDistance:       300,
		Zoom:                1,
		LensRadius:          2.5,
		Sampler:             sampler,
	}

	camera.Initialize()

	var lightSources []render.LightSource
	light := lighting.AreaLight{Object: disk, Material: emmisiveMaterial}
	light.Shadows = true
	// light := render.pointlight{
	// 	location:   geometry.point{100, 250, -150},
	// 	basiclight: render.basiclight{shadows: true, color: render.white, radiancescalingfactor: 2.0}}
	light2 := render.PointLight{
		Location:   geometry.Point{0, 25, -300},
		BasicLight: render.BasicLight{Color: render.WHITE, RadianceScalingFactor: 2.0}}
	light.Initialize()
	light2.Initialize()
	lightSources = append(lightSources, &light)
	lightSources = append(lightSources, &light2)

	world := render.World{Camera: &camera, Lights: lightSources, Objects: objects}

	world.Shading = "area"
	ambientLight := render.AmbientLight{
		BasicLight: render.BasicLight{RadianceScalingFactor: 1.0},
		// MinimumLight: 0.0,
		// Sampler:      sampler,
	}
	ambientLight.Initialize()
	world.AmbientLight = &ambientLight

	render.MultiThreadedMain(&world, samples)

	sceneIo.WritePNG(&camera, output, 0.9)
}
