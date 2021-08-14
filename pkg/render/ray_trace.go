package render

import (
	geometry "basic-ray/pkg/geometry"
	pb "github.com/cheggaaa/pb/v3"
	"math"
)

/*
func Main(origin geometry.Point, lightSources []LightSource, camera *Camera, objects []GeometricObject) {
	var light Photon
	bar := pb.StartNew(len(*camera.Pixels) * len((*camera.Pixels)[0]))

	for i, row := range *camera.Pixels {
		for j, _ := range row {
			ray := geometry.Ray{
				Origin: origin,
				Vector: geometry.Normalize(geometry.CreateVector(GetPoint(camera, i, j), origin)),
			}
			light = Trace(&ray, objects, lightSources, 0)
			(*camera.Pixels)[i][j] = light.rgb // GetWeightedColor()
			bar.Increment()
		}
	}
	bar.Finish()
}
*/

func MultiThreadedMain(world *World, samples int) {
	progress := make(chan bool, 20)
	totalCount := len(*world.Camera.GetPixels()) * len((*world.Camera.GetPixels())[0])
	bar := pb.StartNew(totalCount)
	go renderPixels(world, samples, progress)
	reportProgress(bar, progress, totalCount)
	close(progress)
	bar.Finish()

}

func reportProgress(bar *pb.ProgressBar, progress <-chan bool, totalCount int) {
	for i := 0; i < totalCount; i++ {
		<-progress
		bar.Increment()
	}
}

func renderPixels(world *World, samples int, progress chan<- bool) {
	jobs := make(chan bool, 5)
	viewPlane := world.Camera.GetViewPlane()
	for i := 0; i < viewPlane.HorizontalResolution; i++ {
		for j := 0; j < viewPlane.VerticalResolution; j++ {
			jobs <- true
			rays := world.Camera.GetRays(i, j, samples)
			go renderPixel(rays, world, i, j, progress, jobs)
		}
	}
}

func renderPixel(rays []*geometry.Ray, world *World, i, j int, progress chan<- bool, jobs <-chan bool) {
	var secondaryWeight float64
	secondaryWeight = 1.0 / (float64)(len(rays))
	finalColor := Color{0, 0, 0}
	for _, ray := range rays {
		color := Trace(ray, world, 0)
		finalColor[0] += color[0] * secondaryWeight
		finalColor[1] += color[1] * secondaryWeight
		finalColor[2] += color[2] * secondaryWeight

	}
	// light := Trace(ray, objects, lightSources, 0)
	world.Camera.SetPixel(i, j, finalColor) // GetWeightedColor()
	// complete job
	<-jobs
	// report progress
	progress <- true
}

func Trace(ray *geometry.Ray, world *World, depth int) Color {
	// receiveVector := geometry.Normalize(geometry.ScalarProduct(ray.Vector, -1))
	color := Color{}
	var t float64
	shadeRec := ShadeRec{World: world}
	tmin := float64(math.Inf(1))
	t = float64(math.Inf(1))

	for _, object := range world.Objects {
		intersects := object.Hit(ray, &t, &shadeRec)
		if !intersects || t > tmin {
			continue
		}
		tmin = t
		shadeRec.ObjectHit = true
		shadeRec.HitPoint = geometry.Translate(ray.Origin, geometry.ScalarProduct(ray.Vector, t))
		shadeRec.Material = object.GetMaterial()
	}

	if shadeRec.ObjectHit {
		return shadeRec.Material.Shade(&shadeRec)
	}
	return color
}
