package render

import (
	geometry "basic-ray/pkg/geometry"
	"fmt"
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

func MultiThreadedMain(camera Camera, lightSources []LightSource, objects []GeometricObject, samples int) {
	progress := make(chan bool, 20)
	totalCount := len(*camera.GetPixels()) * len((*camera.GetPixels())[0])
	bar := pb.StartNew(totalCount)
	go renderPixels(objects, lightSources, samples, camera, progress)
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

func renderPixels(objects []GeometricObject, lightSources []LightSource, samples int, camera Camera, progress chan<- bool) {
	jobs := make(chan bool, 5)
	viewPlane := camera.GetViewPlane()
	for i := 0; i < viewPlane.HorizontalResolution; i++ {
		for j := 0; j < viewPlane.VerticalResolution; j++ {
			jobs <- true
			rays := camera.GetRays(i, j, samples)
			go renderPixel(rays, objects, lightSources, camera, i, j, progress, jobs)
		}
	}
}

func renderPixel(rays []geometry.Ray, objects []GeometricObject, lightSources []LightSource, camera Camera, i, j int, progress chan<- bool, jobs <-chan bool) {
	// var secondaryWeight float64 = 0.5 / (float64)(len(rays)-1)
	// var weight float64
	// finalColor := Color{0, 0, 0}
	// for i, ray := range rays {
	// 	light := Trace(ray, objects, lightSources, 0)
	// 	if (light.rgb != Color{0, 0, 0}) {
	// 		fmt.Println(light.rgb)
	// 	}
	// 	if i > 0 {
	// 		finalColor[0] += light.rgb[0] * secondaryWeight
	// 		finalColor[1] += light.rgb[1] * secondaryWeight
	// 		finalColor[2] += light.rgb[2] * secondaryWeight
	// 	} else {
	// 		weight = 0.5
	// 		finalColor[0] += light.rgb[0] * weight
	// 		finalColor[1] += light.rgb[1] * weight
	// 		finalColor[2] += light.rgb[2] * weight
	// 	}

	// }
	light := Trace(&rays[0], objects, lightSources, 0)
	camera.SetPixel(i, j, light.rgb) // GetWeightedColor()
	// complete job
	<-jobs
	// report progress
	progress <- true
}

func Trace(ray *geometry.Ray, objects []GeometricObject, lightSources []LightSource, depth int) Photon {
	receiveVector := geometry.Normalize(geometry.ScalarProduct(ray.Vector, -1))
	photon := Photon{vector: receiveVector}
	var t float64
	shadeRec := ShadeRec{}
	tmin := float64(math.Inf(1))

	for _, object := range objects {
		intersects := object.Hit(ray, &t, &shadeRec)
		if !intersects || t > tmin {
			continue
		}
		tmin = t
		shadeRec.ObjectHit = true
		shadeRec.RGBColor = object.GetColor()
	}

	if shadeRec.ObjectHit {
		fmt.Println("Hit")
		photon.rgb = shadeRec.RGBColor
	}
	return photon
}
